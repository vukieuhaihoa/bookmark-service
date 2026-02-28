package queue

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
)

func TestRepository_PushMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectedError error
		verifyFunc    func(ctx context.Context, redisClient *redis.Client)
	}{
		{
			name: "successful push message",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				return redisClient
			},

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				result, err := redisClient.RPop(ctx, "test_queue").Result()

				assert.Nil(t, err)

				assert.Equal(t, "test_message", result)
			},
		},
		{
			name: "failed push message due to closed Redis client",

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := redisPkg.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRedisClient := tc.setupMock(ctx)
			repo := NewRedisQueue(mockRedisClient, "test_queue")

			err := repo.PushMessage(ctx, []byte("test_message"))

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.Nil(t, err)

			tc.verifyFunc(ctx, mockRedisClient)
		})
	}
}
