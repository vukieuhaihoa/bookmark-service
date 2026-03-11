package link

import (
	"context"
	"errors"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/model"
)

// GetURL retrieves the original URL associated with the given shortened URL code.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations
//   - urlCode: The shortened URL code
//
// Returns:
//   - string: The original URL associated with the shortened code
//   - error: An error object if the retrieval operation fails, otherwise nil
func (l *linkService) GetURL(ctx context.Context, urlCode string) (string, error) {
	s := newrelic.FromContext(ctx).StartSegment("Service_GetURL")
	defer s.End()

	switch {
	case strings.HasPrefix(urlCode, model.RedisShortenPrefix) && len(urlCode) == defaultURLCodeLength+1:
		url, err := l.repo.GetURL(ctx, urlCode)
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}

		return url, err
	case strings.HasPrefix(urlCode, model.BookmarkShortenPrefix):
		bookmark, err := l.bookmarkRepo.GetBookmarkByCodeShortenEncoded(ctx, urlCode)
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	default:
		return "", ErrCodeNotFound
	}
}
