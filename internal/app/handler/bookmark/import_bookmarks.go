package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/csv"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/service/queue"
)

const LimitCSVFileSize = 10 << 20 // 10 MiB

// ImportBookmarks handles the HTTP request to import bookmarks from a file for the authenticated user.
// @Summary      Import bookmarks from a file
// @Description  Import bookmarks from a file for the authenticated user
// @Tags         Bookmarks
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Bookmarks file to import"
// @Success      200   {object}	object{message=string}
// @Failure      400   {object}  object{message=string}
// @Failure      401   {object}  object{message=string}
// @Failure      500   {object}  object{message=string}
// @Security		 Bearer
// @Router       /v1/bookmarks/import [post]
func (h *bookmarkHandler) ImportBookmarks(c *gin.Context) {
	// get user ID from JWT claims
	uid, err := utils.GetUserIDFromJWTClaims(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user ID from JWT claims")
		c.JSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		return
	}

	// get file from form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Message{
			Message: "invalid file.",
		})
		return
	}

	// check file size
	if file.Size > LimitCSVFileSize {
		c.JSON(http.StatusBadRequest, common.Message{
			Message: "file size exceeds the limit.",
		})
		return
	}

	// parse CSV file and convert to import input
	var importInput []*queue.ImportBookmarkInput
	err = csv.ParseFromMultipartFile(file, &importInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse CSV file")
		c.JSON(http.StatusBadRequest, common.Message{
			Message: "failed to parse CSV file.",
		})
		return
	}
	err = h.validator.Var(importInput, "dive")
	if err != nil {
		log.Error().Err(err).Msg("Validation failed for import input")
		c.JSON(http.StatusBadRequest, common.InputFieldError(err))
		return
	}

	err = h.queueSvc.SendImportBookmarkJob(c, uid, importInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to push import bookmark message to queue")
		c.JSON(http.StatusInternalServerError, common.Message{
			Message: "Failed to process bookmark imports.",
		})
		return
	}

	c.JSON(http.StatusOK, common.Message{
		Message: "Successfully sent bookmark imports to queue!",
	})
}
