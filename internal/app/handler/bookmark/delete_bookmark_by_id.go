package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
)

// DeleteBookmarkByID handles the HTTP request to delete a bookmark by its ID for the authenticated user.
// @Summary      Delete bookmark by ID
// @Description  Delete bookmark by ID for authenticated user
// @Tags         Bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Bookmark ID"
// @Success      200  {object}  object{message=string}
// @Failure      400  {object}  object{message=string}
// @Failure      401  {object}  object{message=string}
// @Failure      500  {object}  object{message=string}
// @Security     Bearer
// @Router       /v1/bookmarks/{id} [delete]
func (b *bookmarkHandler) DeleteBookmarkByID(c *gin.Context) {
	s := newrelic.FromContext(c).StartSegment("Handler_DeleteBookmarkByID")
	defer s.End()

	bookmarkID := c.Param("id")
	if bookmarkID == "" {
		c.JSON(http.StatusBadRequest, common.Message{
			Message: "Bookmark ID is required",
		})
		return
	}

	userID, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	err = b.svc.DeleteBookmarkByID(c, bookmarkID, userID)
	switch {
	case errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusBadRequest, common.InputErrorResponse)
		return
	case err == nil:
	default:
		log.Error().Err(err).Msg("Failed to delete bookmark by ID")
		c.JSON(http.StatusInternalServerError, common.InternalErrorResponse)
		return
	}

	c.JSON(http.StatusOK, common.Message{
		Message: "Success",
	})
}
