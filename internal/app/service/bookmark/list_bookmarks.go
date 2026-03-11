package bookmark

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/model"
)

// ListBookmarks retrieves a list of bookmarks for a specific user based on the provided query options.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - userID: The ID of the user whose bookmarks are to be retrieved.
//   - opts: Query options for filtering, sorting, and pagination.
//
// Returns:
//   - []*model.Bookmark: A slice of bookmark models.
//   - error: An error if the retrieval fails.
func (b *bookmarkService) ListBookmarks(ctx context.Context, userID string, opts *common.QueryOptions) ([]*model.Bookmark, error) {
	s := newrelic.FromContext(ctx).StartSegment("Service_ListBookmarks")
	defer s.End()

	return b.repo.ListBookmarks(ctx, userID, opts)
}
