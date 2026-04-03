package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/middleware"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

type NewsHandler struct {
	newsService *service.NewsService
}

func NewNewsHandler(newsService *service.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

func (h *NewsHandler) GetHeadlines(c *gin.Context) {
	headlines, err := h.newsService.GetHeadlines(c.Request.Context())
	if err != nil {
		response.InternalError(c, "failed to get headlines")
		return
	}
	if headlines == nil {
		response.Success(c, []interface{}{})
		return
	}
	response.Success(c, headlines)
}

func (h *NewsHandler) GetHeadlineHistory(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		response.BadRequest(c, "date parameter is required (YYYY-MM-DD)")
		return
	}

	headlines, err := h.newsService.GetHeadlineHistory(c.Request.Context(), date)
	if err != nil {
		response.InternalError(c, "failed to get headline history")
		return
	}
	response.Success(c, headlines)
}

func (h *NewsHandler) ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	category := c.Query("category")
	section := c.Query("section")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := h.newsService.List(c.Request.Context(), page, size, category, section)
	if err != nil {
		response.InternalError(c, "failed to list news")
		return
	}

	response.SuccessPage(c, list, total, page, size)
}

func (h *NewsHandler) GetNews(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	news, err := h.newsService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "failed to get news")
		return
	}
	if news == nil {
		response.NotFound(c, "news not found")
		return
	}

	response.Success(c, news)
}

func (h *NewsHandler) SearchNews(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.BadRequest(c, "search query is required")
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

	list, total, err := h.newsService.Search(c.Request.Context(), q, page, size)
	if err != nil {
		response.InternalError(c, "search failed")
		return
	}

	response.SuccessPage(c, list, total, page, size)
}

func (h *NewsHandler) MarkRead(c *gin.Context) {
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

	if err := h.newsService.MarkRead(c.Request.Context(), userID, newsID); err != nil {
		response.InternalError(c, "failed to mark read")
		return
	}

	response.Success(c, gin.H{"marked": true})
}
