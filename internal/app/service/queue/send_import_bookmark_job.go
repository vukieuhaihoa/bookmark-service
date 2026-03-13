package queue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/array"
)

const DefaultImportBookmarkBatchSize = 5

var (
	ErrMarshalMessage = errors.New("failed to marshal message")
)

// ImportMessage represents the structure of the message that will be sent to the queue for importing bookmarks. It contains the user ID (UID) and a slice of ImportBookmarkInput, which holds the details of each bookmark to be imported. This struct is used to encapsulate the data that will be processed by the worker responsible for handling bookmark imports from the queue.
type ImportMessage struct {
	UID       string                 `json:"user_id"`
	Bookmarks []*ImportBookmarkInput `json:"bookmarks"`
}

// ImportBookmarkInput represents the structure of a bookmark to be imported. It includes a description and a URL, both of which are required fields with specific validation rules. The description must be a string with a maximum length of 255 characters, while the URL must be a valid URL string with a maximum length of 1024 characters. This struct is used to capture the details of each bookmark that is being imported through the API.
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
func (q *queueService) SendImportBookmarkJob(ctx context.Context, uid string, bookmarks []*ImportBookmarkInput) error {
	s := newrelic.FromContext(ctx).StartSegment("Service_SendImportBookmarkJob")
	defer s.End()

	batches := array.SplitIntoBatches(bookmarks, DefaultImportBookmarkBatchSize)
	for _, batch := range batches {
		if err := q.sendJob(ctx, uid, batch); err != nil {
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
		return ErrMarshalMessage
	}

	return s.repo.PushMessage(ctx, messageBytes)
}
