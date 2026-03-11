package bookmark

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/model"
)

// CreateBookmark creates a new bookmark with the provided information.
// Parameters:
//   - ctx: The context for managing request-scoped values and cancellation.
//   - url: The URL of the bookmark.
//   - description: A description of the bookmark.
//   - userID: The ID of the user who owns the bookmark.
//
// Returns:
//   - *model.Bookmark: The created bookmark model.
//   - error: An error if the creation fails, otherwise nil.
func (b *bookmarkService) CreateBookmark(ctx context.Context, url, description, userID string) (*model.Bookmark, error) {
	s := newrelic.FromContext(ctx).StartSegment("Service_CreateBookmark")
	defer s.End()

	newBookmark := &model.Bookmark{
		URL:         url,
		Description: description,
		UserID:      userID,
	}

	createdBookmark, err := b.repo.CreateBookmark(ctx, newBookmark)
	if err != nil {
		return nil, err
	}

	return createdBookmark, nil
}
