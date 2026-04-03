package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/middleware"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type sendSMSRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *AuthHandler) SendSMS(c *gin.Context) {
	var req sendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "phone is required")
		return
	}

	if len(req.Phone) < 10 {
		response.BadRequest(c, "invalid phone number")
		return
	}

	if err := h.authService.SendSMS(c.Request.Context(), req.Phone); err != nil {
		response.InternalError(c, "failed to send SMS")
		return
	}

	response.Success(c, gin.H{"sent": true})
}

type smsLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (h *AuthHandler) LoginWithSMS(c *gin.Context) {
	var req smsLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "phone and code are required")
		return
	}

	token, user, err := h.authService.LoginWithSMS(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		response.Error(c, 401, response.CodeAuthInvalid, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

type refreshRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "token is required")
		return
	}

	newToken, err := h.authService.RefreshToken(c.Request.Context(), req.Token)
	if err != nil {
		response.Error(c, 401, response.CodeAuthExpired, "invalid or expired token")
		return
	}

	response.Success(c, gin.H{"token": newToken})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil || user == nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}
