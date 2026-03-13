package bookmark

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/model"
)

// GetBookmarkByCodeShortenEncoded retrieves a bookmark by its unique code_shorten_encoded.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - code: The unique code_shorten_encoded of the bookmark to retrieve.
//
// Returns:
//   - A pointer to the Bookmark model if found.
//   - An error if the bookmark is not found or if a database error occurs.
func (b *bookmarkRepository) GetBookmarkByCodeShortenEncoded(ctx context.Context, code string) (*model.Bookmark, error) {
	s := newrelic.FromContext(ctx).StartSegment("Repo_GetBookmarkByCodeShortenEncoded")
	defer s.End()

	res := &model.Bookmark{}
	err := b.db.WithContext(ctx).Where("code_shorten_encoded = ?", code).First(res).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}

	return res, nil
}
