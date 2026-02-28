package queue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/array"
)

const DefaultImportBookmarkBatchSize = 5

var (
	ErrUnmarshalMessage = errors.New("failed to unmarshal message")
)

type ImportMessage struct {
	UID       string                 `json:"user_id"`
	Bookmarks []*ImportBookmarkInput `json:"bookmarks"`
}

type ImportBookmarkInput struct {
	Description string `csv:"description" binding:"required,lte=255" json:"description"`
	URL         string `csv:"url" binding:"required,url,lte=1024" json:"url"`
}

// SendImportBookmarkJob sends a job to the message queue for importing bookmarks. It takes the user ID and a slice of ImportBookmarkInput, splits the bookmarks into batches, and sends each batch as a separate job to the queue. If any step fails, it returns an appropriate error.
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation.
//   - uid: The unique identifier of the user for whom the bookmarks are being imported.
//   - bookmarks: A slice of ImportBookmarkInput containing the details of the bookmarks to be imported.
//
// Returns an error if the operation fails, otherwise nil.
func (s *queueService) SendImportBookmarkJob(ctx context.Context, uid string, bookmarks []*ImportBookmarkInput) error {
	batches := array.SplitIntoBatches(bookmarks, DefaultImportBookmarkBatchSize)
	for _, batch := range batches {
		if err := s.sendJob(ctx, uid, batch); err != nil {
			return err
		}
	}
	return nil
}

// sendJob sends a single job to the message queue for importing bookmarks. It takes the user ID and a batch of bookmarks, constructs an ImportMessage, marshals it to JSON, and pushes it to the queue. If any step fails, it returns an appropriate error.
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation.
//   - uid: The unique identifier of the user for whom the bookmarks are being imported.
//   - bookmarks: A slice of ImportBookmarkInput containing the details of the bookmarks to be imported.
//
// Returns an error if the operation fails, otherwise nil.
func (s *queueService) sendJob(ctx context.Context, uid string, bookmarks []*ImportBookmarkInput) error {
	message := ImportMessage{
		UID:       uid,
		Bookmarks: bookmarks,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return ErrUnmarshalMessage
	}

	return s.repo.PushMessage(ctx, messageBytes)
}
