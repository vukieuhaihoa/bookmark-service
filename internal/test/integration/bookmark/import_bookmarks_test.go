package bookmark

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	middleware "github.com/vukieuhaihoa/bookmark-libs/middlewares"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils/mocks"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-service/internal/api"
	"github.com/vukieuhaihoa/bookmark-service/pkg/testutils"
)

func TestBookmarkEndpoint_ImportBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		fileContent string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client

		setupTestHTTP func(api api.Engine, body *bytes.Buffer, writer *multipart.Writer) *httptest.ResponseRecorder

		setupMockJWTValidator func(t *testing.T) *mocks.JWTValidator

		expectedStatusCode      int
		expectedMessageResponse string

		verifyRedisQueue func(ctx context.Context, mock *redis.Client)
	}{
		{
			name: "successful import bookmarks",

			fileContent: "description,url\nExample Website,https://example.com",

			setupTestHTTP: func(api api.Engine, body *bytes.Buffer, writer *multipart.Writer) *httptest.ResponseRecorder {
				req := httptest.NewRequest("POST", "/v1/bookmarks/import", body)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				req.Header.Set("Content-Type", writer.FormDataContentType())
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusOK,
			expectedMessageResponse: `"message":"Successfully sent bookmark imports to queue!"`,

			verifyRedisQueue: func(ctx context.Context, mock *redis.Client) {
				assert.Equal(t, int64(1), mock.LLen(ctx, "bookmark_import_queue").Val())
				assert.Equal(t, `{"user_id":"4d9326d6-980c-4c62-9709-dbc70a82cbfe","bookmarks":[{"description":"Example Website","url":"https://example.com"}]}`, mock.RPop(ctx, "bookmark_import_queue").Val())
			},
		},
		{
			name: "failed import bookmarks - empty csv file",

			fileContent: "",

			setupTestHTTP: func(api api.Engine, body *bytes.Buffer, writer *multipart.Writer) *httptest.ResponseRecorder {
				req := httptest.NewRequest("POST", "/v1/bookmarks/import", body)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				req.Header.Set("Content-Type", writer.FormDataContentType())
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: `"message":"failed to parse CSV file."`,
		},
		{
			name: "import bookmark failed - invalid token",
			setupTestHTTP: func(api api.Engine, body *bytes.Buffer, writer *multipart.Writer) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/bookmarks/import", body)
				req.Header.Set("Authorization", "Bearer invalid_jwt_token")
				req.Header.Set("Content-Type", writer.FormDataContentType())
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "invalid_jwt_token").Return(nil, assert.AnError)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusUnauthorized,
			expectedMessageResponse: `"message":"Invalid token"`,
		},
		{
			name: "import bookmark failed - token does not contain user ID",

			setupTestHTTP: func(api api.Engine, body *bytes.Buffer, writer *multipart.Writer) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/bookmarks/import", body)
				req.Header.Set("Authorization", "Bearer token_without_user_id")
				req.Header.Set("Content-Type", writer.FormDataContentType())
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "token_without_user_id").Return(jwt.MapClaims{}, nil)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusUnauthorized,
			expectedMessageResponse: `"message":"Unauthorized"`,
		},
		{
			name: "import bookmark failed - rate limit exceeded",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				key := fmt.Sprintf(middleware.RateLimitKeyFormat, "4d9326d6-980c-4c62-9709-dbc70a82cbfe")
				redisClient.Set(ctx, key, middleware.UserIDRateLimitMaxCount, middleware.UserIDRateLimitInterval)
				return redisClient
			},

			setupTestHTTP: func(api api.Engine, body *bytes.Buffer, writer *multipart.Writer) *httptest.ResponseRecorder {
				// Setup HTTP request and recorder
				req := httptest.NewRequest("POST", "/v1/bookmarks/import", body)
				req.Header.Set("Authorization", "Bearer valid_jwt_token")
				req.Header.Set("Content-Type", writer.FormDataContentType())
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			setupMockJWTValidator: func(t *testing.T) *mocks.JWTValidator {
				jwtValidator := mocks.NewJWTValidator(t)
				jwtValidator.On("ValidateToken", "valid_jwt_token").Return(jwt.MapClaims{"sub": "4d9326d6-980c-4c62-9709-dbc70a82cbfe"}, nil)
				return jwtValidator
			},

			expectedStatusCode:      http.StatusTooManyRequests,
			expectedMessageResponse: `"error":"Too many requests. Please try again later."`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			writer, body, tmpfile := testutils.CreateMultipartRequest(t, tc.fileContent)
			defer os.Remove(tmpfile)

			jwtValidator := tc.setupMockJWTValidator(t)
			redisClient := redisPkg.InitMockRedis(t)
			if tc.setupMockRedis != nil {
				redisClient = tc.setupMockRedis(ctx, redisClient)
			}

			apiEngine := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg: &api.Config{
					ServiceName: "bookmark_service",
					InstanceID:  "test_instance_id_1",
				},
				RedisClient:     redisClient,
				SqlDB:           nil,
				RandomCodeGen:   nil,
				PasswordHashing: nil,
				JWTGenerator:    nil,
				JWTValidator:    jwtValidator,
			})

			respRec := tc.setupTestHTTP(apiEngine, body, writer)

			assert.Equal(t, tc.expectedStatusCode, respRec.Code)
			assert.Contains(t, respRec.Body.String(), tc.expectedMessageResponse)

			if tc.verifyRedisQueue != nil {
				tc.verifyRedisQueue(ctx, redisClient)
			}
		})
	}
}
