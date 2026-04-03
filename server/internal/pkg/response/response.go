package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard error codes.
const (
	CodeSuccess         = 0
	CodeAuthRequired    = 1001
	CodeAuthInvalid     = 1002
	CodeAuthExpired     = 1003
	CodeForbidden       = 1004
	CodeParamInvalid    = 2001
	CodeParamMissing    = 2002
	CodeBusinessError   = 3001
	CodeNotFound        = 3002
	CodeDuplicate       = 3003
	CodeInsufficientBal = 3004
	CodeInternalError   = 5001
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessPage(c *gin.Context, list interface{}, total, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PageData{
			List:  list,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}

func Error(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: msg,
	})
}

func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, CodeParamInvalid, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, CodeAuthRequired, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, CodeForbidden, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, CodeNotFound, msg)
}

func InternalError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, msg)
}
