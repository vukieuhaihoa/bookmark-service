package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mockQueueRepo "github.com/vukieuhaihoa/bookmark-service/internal/app/repository/queue/mocks"
)

func TestService_SendImportBookmarkJob(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		uid       string
		bookmarks []*ImportBookmarkInput

		setupMock func(ctx context.Context) *mockQueueRepo.Repository

		expectedError error
	}{
		{
			name: "successful send import bookmark job - single batch",

			uid: "user123",
			bookmarks: []*ImportBookmarkInput{
				{
					URL:         "https://example.com",
					Description: "An example bookmark",
				},
			},

			setupMock: func(ctx context.Context) *mockQueueRepo.Repository {
				repoMock := mockQueueRepo.NewRepository(t)
				input := ImportMessage{
					UID: "user123",
					Bookmarks: []*ImportBookmarkInput{
						{
							URL:         "https://example.com",
							Description: "An example bookmark",
						},
					},
				}
				inputBytes, err := json.Marshal(input)
				assert.NoError(t, err)

				repoMock.On("PushMessage", ctx, inputBytes).Return(nil)
				return repoMock
			},

			expectedError: nil,
		},
		{
			name: "successful send import bookmark job - multiple batches",

			uid: "user123",
			bookmarks: []*ImportBookmarkInput{
				{
					URL:         "https://example.com/1",
					Description: "Bookmark 1",
				},
				{
					URL:         "https://example.com/2",
					Description: "Bookmark 2",
				},
				{
					URL:         "https://example.com/3",
					Description: "Bookmark 3",
				},
				{
					URL:         "https://example.com/4",
					Description: "Bookmark 4",
				},
				{
					URL:         "https://example.com/5",
					Description: "Bookmark 5",
				},
				{
					URL:         "https://example.com/6",
					Description: "Bookmark 6",
				},
			},

			setupMock: func(ctx context.Context) *mockQueueRepo.Repository {
				repoMock := mockQueueRepo.NewRepository(t)
				input1 := ImportMessage{
					UID: "user123",
					Bookmarks: []*ImportBookmarkInput{
						{
							URL:         "https://example.com/1",
							Description: "Bookmark 1",
						},
						{
							URL:         "https://example.com/2",
							Description: "Bookmark 2",
						},
						{
							URL:         "https://example.com/3",
							Description: "Bookmark 3",
						},
						{
							URL:         "https://example.com/4",
							Description: "Bookmark 4",
						},
						{
							URL:         "https://example.com/5",
							Description: "Bookmark 5",
						},
					},
				}

				inputBytes1, err := json.Marshal(input1)
				assert.NoError(t, err)

				input2 := ImportMessage{
					UID: "user123",
					Bookmarks: []*ImportBookmarkInput{
						{
							URL:         "https://example.com/6",
							Description: "Bookmark 6",
						},
					},
				}

				inputBytes2, err := json.Marshal(input2)
				assert.NoError(t, err)

				repoMock.On("PushMessage", ctx, inputBytes1).Return(nil).Once()
				repoMock.On("PushMessage", ctx, inputBytes2).Return(nil).Once()

				return repoMock
			},

			expectedError: nil,
		},
		{
			name: "failed to send import bookmark job due to repository error",
			uid:  "user123",
			bookmarks: []*ImportBookmarkInput{
				{
					URL:         "https://example.com",
					Description: "An example bookmark",
				},
			},

			setupMock: func(ctx context.Context) *mockQueueRepo.Repository {
				repoMock := mockQueueRepo.NewRepository(t)
				input := ImportMessage{
					UID: "user123",
					Bookmarks: []*ImportBookmarkInput{
						{
							URL:         "https://example.com",
							Description: "An example bookmark",
						},
					},
				}
				inputBytes, err := json.Marshal(input)
				assert.NoError(t, err)

				repoMock.On("PushMessage", ctx, inputBytes).Return(errors.New("repository error"))
				return repoMock
			},

			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMock(ctx)
			service := NewQueueService(mockRepo)

			err := service.SendImportBookmarkJob(ctx, tc.uid, tc.bookmarks)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				return
			}

			assert.Nil(t, err)
		})
	}
}
