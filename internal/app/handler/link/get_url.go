package link

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	service "github.com/vukieuhaihoa/bookmark-service/internal/app/service/link"
)

// GetURL handles the request to retrieve the original URL from a shortened code.
// It takes a Gin context as input and processes the request to fetch the original URL.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response
//
// @Summary Get original URL
// @Description Retrieves the original URL from a given shortened code.
// @Tags URL
// @Accept json
// @Produce json
// @Param code path string true "Shortened code"
// @Success 200 {object} shortenURLResponse
// @Failure 400 {object} shortenURLResponse
// @Failure 500 {object} shortenURLResponse
// @Router /v1/links/redirect/{code} [get]
func (h *linkHandler) GetURL(c *gin.Context) {
	s := newrelic.FromContext(c).StartSegment("Handler_GetURL")
	defer s.End()

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: InValidRequestPayload.Error(),
		})
		return
	}

	url, err := h.linkSvc.GetURL(c, code)
	// Handle different error cases with switch statement
	switch {
	case errors.Is(err, service.ErrCodeNotFound) || errors.Is(err, dbutils.ErrRecordNotFoundType):
		c.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: ErrCodeNotFound.Error(),
		})
		return
	case err == nil: // MUST: to redirect when no error
	default:
		log.Error().Str("code", code).Err(err).Msg("service return error when get original url")
		c.JSON(http.StatusInternalServerError, shortenURLResponse{
			Message: InternalServerError.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, url)
}
