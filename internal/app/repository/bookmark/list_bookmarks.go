package bookmark

import (
	"context"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/model"
)

// ListBookmarks retrieves a list of bookmarks from the database for a specific user.
// It takes a context, a user ID, and query options for pagination and sorting.
// Returns a slice of bookmark models and an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - userID: The ID of the user whose bookmarks are to be retrieved.
//   - opts: QueryOptions containing pagination and sorting details.
//
// Returns:
//   - []*model.Bookmark: A slice of bookmark models.
//   - error: An error if the retrieval fails, otherwise nil.
func (r *bookmarkRepository) ListBookmarks(ctx context.Context, userID string, opts *common.QueryOptions) ([]*model.Bookmark, error) {
	s := newrelic.FromContext(ctx).StartSegment("Repo_ListBookmarks")
	defer s.End()

	var bookmarks []*model.Bookmark

	query := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID)

	for _, sortField := range opts.Sorting {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}

	if err := query.Count(&opts.Total).Error; err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	if err := query.Offset((opts.Page - 1) * opts.Limit).Limit(opts.Limit).Find(&bookmarks).Error; err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return bookmarks, nil
}
