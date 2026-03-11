package queue

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// PushMessage pushes a message to the Redis queue.
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation.
//   - message: The message to be pushed to the queue, represented as a byte slice.
//
// Returns an error if the operation fails.
func (r *redisQueue) PushMessage(ctx context.Context, message []byte) error {
	nrs := newrelic.FromContext(ctx).StartSegment("Repo_PushMessage")
	defer nrs.End()

	return r.c.LPush(ctx, r.queueName, message).Err()
}
