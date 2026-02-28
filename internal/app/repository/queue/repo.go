package queue

import "context"

// Repository defines the interface for a message queue repository.
// It provides methods to interact with the message queue, such as pushing messages.
//
//go:generate mockery --name=Repository --output=./mocks --filename=repo.go
type Repository interface {
	// PushMessage pushes a message to the queue.
	// Parameters:
	//   - ctx: The context for managing request deadlines and cancellation.
	//   - message: The message to be pushed to the queue, represented as a byte slice.
	// Returns an error if the operation fails.
	PushMessage(ctx context.Context, message []byte) error
}
