package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/model"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

type AdminHandler struct {
	auditService *service.AuditService
	newsService  *service.NewsService
	eventService *service.EventService
}

func NewAdminHandler(auditService *service.AuditService, newsService *service.NewsService, eventService *service.EventService) *AdminHandler {
	return &AdminHandler{
		auditService: auditService,
		newsService:  newsService,
		eventService: eventService,
	}
}

func (h *AdminHandler) GetAuditQueue(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}

	list, total, err := h.auditService.GetPendingQueue(c.Request.Context(), page, size)
	if err != nil {
		response.InternalError(c, "failed to get audit queue")
		return
	}

	response.SuccessPage(c, list, total, page, size)
}

func (h *AdminHandler) ApproveComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid comment id")
		return
	}

	if err := h.auditService.Approve(c.Request.Context(), id); err != nil {
		response.InternalError(c, "failed to approve")
		return
	}

	response.Success(c, gin.H{"approved": true})
}

type rejectRequest struct {
	Reason string `json:"reason"`
}

func (h *AdminHandler) RejectComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid comment id")
		return
	}

	var req rejectRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.auditService.Reject(c.Request.Context(), id, req.Reason); err != nil {
		response.InternalError(c, "failed to reject")
		return
	}

	response.Success(c, gin.H{"rejected": true})
}

type setHeadlinesRequest struct {
	Headlines []headlineItem `json:"headlines" binding:"required"`
}

type headlineItem struct {
	NewsID   int64  `json:"news_id" binding:"required"`
	Rank     int    `json:"rank"`
	Title    string `json:"title"`
	ExpireAt string `json:"expire_at"` // RFC3339
}

func (h *AdminHandler) SetHeadlines(c *gin.Context) {
	var req setHeadlinesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	var headlines []*model.Headline
	now := time.Now().UTC()
	for _, item := range req.Headlines {
		expireAt := now.Add(24 * time.Hour)
		if item.ExpireAt != "" {
			if t, err := time.Parse(time.RFC3339, item.ExpireAt); err == nil {
				expireAt = t
			}
		}
		headlines = append(headlines, &model.Headline{
			NewsID:   item.NewsID,
			Rank:     item.Rank,
			Title:    item.Title,
			ActiveAt: now,
			ExpireAt: expireAt,
		})
	}

	if err := h.newsService.SetHeadlines(c.Request.Context(), headlines); err != nil {
		response.InternalError(c, "failed to set headlines")
		return
	}

	response.Success(c, gin.H{"set": true})
}

type createEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CoverImage  string `json:"cover_image"`
}

func (h *AdminHandler) CreateEvent(c *gin.Context) {
	var req createEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "title is required")
		return
	}

	event := &model.Event{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		CoverImage:  req.CoverImage,
	}

	id, err := h.eventService.Create(c.Request.Context(), event)
	if err != nil {
		response.InternalError(c, "failed to create event")
		return
	}

	response.Success(c, gin.H{"id": id})
}

type createNodeRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	NodeTime string `json:"node_time" binding:"required"` // RFC3339
	Source   string `json:"source"`
	Veracity int    `json:"veracity"` // 0=待核实, 1=已证实, 2=已辟谣
}

func (h *AdminHandler) CreateEventNode(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid event id")
		return
	}

	var req createNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "title and node_time are required")
		return
	}

	nodeTime, err := time.Parse(time.RFC3339, req.NodeTime)
	if err != nil {
		response.BadRequest(c, "invalid node_time format, use RFC3339")
		return
	}

	node := &model.EventNode{
		EventID:  eventID,
		Title:    req.Title,
		Content:  req.Content,
		NodeTime: nodeTime,
		Source:   req.Source,
		Veracity: req.Veracity,
	}

	id, err := h.eventService.AddNode(c.Request.Context(), node)
	if err != nil {
		response.InternalError(c, "failed to create node")
		return
	}

	response.Success(c, gin.H{"id": id})
}
