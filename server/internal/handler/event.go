package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := h.eventService.List(c.Request.Context(), page, size)
	if err != nil {
		response.InternalError(c, "failed to list events")
		return
	}

	response.SuccessPage(c, list, total, page, size)
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid event id")
		return
	}

	detail, err := h.eventService.GetWithTimeline(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "failed to get event")
		return
	}
	if detail == nil {
		response.NotFound(c, "event not found")
		return
	}

	response.Success(c, detail)
}
