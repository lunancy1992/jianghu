package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/middleware"
	"github.com/lunancy1992/jianghu-server/internal/model"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

type EvidenceHandler struct {
	eventService *service.EventService
	dataDir      string
}

func NewEvidenceHandler(eventService *service.EventService, dataDir string) *EvidenceHandler {
	return &EvidenceHandler{eventService: eventService, dataDir: dataDir}
}

// Upload handles multipart file upload and saves to local disk.
func (h *EvidenceHandler) Upload(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	// Create upload directory
	uploadDir := filepath.Join(h.dataDir, "uploads", "evidence")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "failed to create upload directory")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	url := "/uploads/evidence/" + filename

	response.Success(c, gin.H{
		"url":      url,
		"filename": filename,
	})
}

type addEvidenceRequest struct {
	Type        string `json:"type"` // image, document, link
	URL         string `json:"url" binding:"required"`
	Description string `json:"description"`
}

// AddEvidence attaches evidence to an event.
func (h *EvidenceHandler) AddEvidence(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid event id")
		return
	}

	var req addEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "url is required")
		return
	}

	if req.Type == "" {
		req.Type = "image"
	}

	evidence := &model.Evidence{
		EventID:     eventID,
		UserID:      userID,
		Type:        req.Type,
		URL:         req.URL,
		Description: req.Description,
		Status:      0, // pending
	}

	id, err := h.eventService.AddEvidence(c.Request.Context(), evidence)
	if err != nil {
		response.InternalError(c, "failed to add evidence")
		return
	}
	if id == 0 {
		response.NotFound(c, "event not found")
		return
	}

	response.Success(c, gin.H{"id": id})
}
