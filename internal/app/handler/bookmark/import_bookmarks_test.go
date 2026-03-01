package bookmark

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/vukieuhaihoa/bookmark-service/internal/app/service/queue"
	mockQueueSvc "github.com/vukieuhaihoa/bookmark-service/internal/app/service/queue/mocks"
	"github.com/vukieuhaihoa/bookmark-service/pkg/testutils"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.SetTagName("binding")
}

func TestHandler_ImportBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer)

		setupMockSvc func(ctx *gin.Context) *mockQueueSvc.Service

		fileContent string

		expectedCode     int
		expectedResponse string
	}{
		{
			name: "successful import bookmarks",

			setupRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "user-123"})

				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},

			setupMockSvc: func(ctx *gin.Context) *mockQueueSvc.Service {
				svcMock := mockQueueSvc.NewService(t)
				svcMock.On("SendImportBookmarkJob", ctx, "user-123", []*queue.ImportBookmarkInput{
					{
						URL:         "https://example.com",
						Description: "Example Website",
					},
					{
						URL:         "https://another.com",
						Description: "Another Website",
					},
				}).Return(nil)
				return svcMock
			},

			fileContent: "description,url\nExample Website,https://example.com\nAnother Website,https://another.com",

			expectedCode:     http.StatusOK,
			expectedResponse: `{"message":"Successfully sent bookmark imports to queue!"}`,
		},
		{
			name: "error case - svc returns error",

			setupRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "user-123"})

				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},

			setupMockSvc: func(ctx *gin.Context) *mockQueueSvc.Service {
				svcMock := mockQueueSvc.NewService(t)
				svcMock.On("SendImportBookmarkJob", ctx, "user-123", []*queue.ImportBookmarkInput{
					{
						URL:         "https://example.com",
						Description: "Example Website",
					},
					{
						URL:         "https://another.com",
						Description: "Another Website",
					},
				}).Return(assert.AnError)
				return svcMock
			},

			fileContent: "description,url\nExample Website,https://example.com\nAnother Website,https://another.com",

			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"message":"Failed to process bookmark imports."}`,
		},
		{
			name: "error case - empty CSV file",

			setupRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "user-123"})

				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},

			setupMockSvc: func(ctx *gin.Context) *mockQueueSvc.Service {
				return mockQueueSvc.NewService(t)
			},

			fileContent: "",

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"failed to parse CSV file."}`,
		},
		{
			name: "error case - a row is invalid format",

			setupRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
				// Simulate authenticated user
				ctx.Set("claims", jwt.MapClaims{"sub": "user-123"})

				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},

			setupMockSvc: func(ctx *gin.Context) *mockQueueSvc.Service {
				return mockQueueSvc.NewService(t)
			},

			fileContent: "description,url\nExample Website,invalid-url\nAnother Website,https://another.com",

			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid input fields","details":["URL is invalid (url)"]}`,
		},
		{
			name: "error case - unauthorized user",

			setupRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
				// No authenticated user set

				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},

			setupMockSvc: func(ctx *gin.Context) *mockQueueSvc.Service {
				return mockQueueSvc.NewService(t)
			},

			fileContent: "description,url\nExample Website,https://example.com\nAnother Website,https://another.com",

			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"message":"Unauthorized"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			writer, body := testutils.CreateMultipartRequest(t, tc.fileContent)

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx, body, writer)
			queueSvc := tc.setupMockSvc(ctx)

			handler := &bookmarkHandler{
				queueSvc:  queueSvc,
				validator: validate,
			}

			handler.ImportBookmarks(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
		})
	}
}
