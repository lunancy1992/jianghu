package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/middleware"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"github.com/lunancy1992/jianghu-server/internal/service"
)

type CoinHandler struct {
	coinService *service.CoinService
}

func NewCoinHandler(coinService *service.CoinService) *CoinHandler {
	return &CoinHandler{coinService: coinService}
}

func (h *CoinHandler) GetBalance(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	balance, err := h.coinService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "failed to get balance")
		return
	}

	response.Success(c, gin.H{"balance": balance})
}

func (h *CoinHandler) GetTransactions(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
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

	list, total, err := h.coinService.GetTransactions(c.Request.Context(), userID, page, size)
	if err != nil {
		response.InternalError(c, "failed to get transactions")
		return
	}

	response.SuccessPage(c, list, total, page, size)
}
