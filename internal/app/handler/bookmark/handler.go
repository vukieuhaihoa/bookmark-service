package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/service/bookmark"
	"github.com/vukieuhaihoa/bookmark-service/internal/app/service/queue"
)

// Handler defines the interface for bookmark-related HTTP handlers.
type Handler interface {

	// CreateBookmark handles the HTTP request to create a new bookmark for the authenticated user.
	CreateBookmark(c *gin.Context)

	// ListBookmarks handles the HTTP request to list bookmarks for the authenticated user.
	ListBookmarks(c *gin.Context)

	// UpdateBookmarkByID handles the HTTP request to update a bookmark by its ID for the authenticated user.
	UpdateBookmarkByID(c *gin.Context)

	// DeleteBookmarkByID handles the HTTP request to delete a bookmark by its ID for the authenticated user.
	DeleteBookmarkByID(c *gin.Context)

	// ImportBookmarks handles the HTTP request to import bookmarks from a file for the authenticated user.
	ImportBookmarks(c *gin.Context)
}

// bookmarkHandler is the concrete implementation of the Handler interface.
type bookmarkHandler struct {
	svc       bookmark.Service
	queueSvc  queue.Service
	validator *validator.Validate
}

// NewBookmarkHandler creates a new instance of bookmarkHandler.
// It takes a bookmark.Service and a queue.Service as parameters and returns a Handler interface.
func NewBookmarkHandler(svc bookmark.Service, queueSvc queue.Service, validator *validator.Validate) Handler {
	return &bookmarkHandler{
		svc:       svc,
		queueSvc:  queueSvc,
		validator: validator,
	}
}
