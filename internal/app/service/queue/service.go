package queue

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-service/internal/app/repository/queue"
)

// Service defines the interface for interacting with the message queue.
//
//go:generate mockery --name Service --filename queue_service.go --output ./mocks
type Service interface {
	// SendImportBookmarkJob sends a job to the message queue for importing bookmarks.
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellation.
	//   - uid: The unique identifier of the user for whom the bookmarks are being imported.
	//   - bookmarks: A slice of ImportBookmarkInput containing the details of the bookmarks to be imported.
	//
	// Returns an error if the operation fails, otherwise nil.
	SendImportBookmarkJob(ctx context.Context, uid string, bookmarks []*ImportBookmarkInput) error
}

// queueService is the concrete implementation of the Service interface that interacts with the message queue.
type queueService struct {
	repo queue.Repository
}

// NewQueueService creates a new instance of queueService with the provided queue.Repository.
// Parameters:
//   - repo: An implementation of the queue.Repository interface that will be used by the service.
//
// Returns a Service interface that can be used to interact with the message queue.
func NewQueueService(repo queue.Repository) Service {
	return &queueService{
		repo: repo,
	}
}
