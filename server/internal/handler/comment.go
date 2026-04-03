package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/middleware"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) ListComments(c *gin.Context) {
	newsID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := h.commentService.List(c.Request.Context(), newsID, page, size)
	if err != nil {
		response.InternalError(c, "failed to list comments")
		return
	}

	response.SuccessPage(c, list, total, page, size)
}

type createCommentRequest struct {
	Content string `json:"content" binding:"required"`
	Stance  string `json:"stance"`
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	newsID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "content is required")
		return
	}

	if req.Stance == "" {
		req.Stance = "neutral"
	}

	comment, err := h.commentService.Create(c.Request.Context(), newsID, userID, req.Content, req.Stance)
	if err != nil {
		response.Error(c, 400, response.CodeBusinessError, err.Error())
		return
	}

	response.Success(c, comment)
}

type likeRequest struct {
	CommentID int64 `json:"comment_id" binding:"required"`
}

func (h *CommentHandler) Like(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req likeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "comment_id is required")
		return
	}

	if err := h.commentService.Like(c.Request.Context(), userID, req.CommentID); err != nil {
		response.InternalError(c, "failed to like")
		return
	}

	response.Success(c, gin.H{"liked": true})
}

func (h *CommentHandler) Unlike(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req likeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "comment_id is required")
		return
	}

	if err := h.commentService.Unlike(c.Request.Context(), userID, req.CommentID); err != nil {
		response.InternalError(c, "failed to unlike")
		return
	}

	response.Success(c, gin.H{"unliked": true})
}
