package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSON(c *gin.Context, status int, success bool, message string, data interface{}, errs interface{}) {
	c.JSON(status, Response{
		Success: success,
		Message: message,
		Data:    data,
		Errors:  errs,
	})
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, true, "Success", data, nil)
}

func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, true, "Created", data, nil)
}

func Error(c *gin.Context, status int, message string) {
	JSON(c, status, false, message, nil, nil)
}

func ValidationError(c *gin.Context, errs interface{}) {
	JSON(c, http.StatusBadRequest, false, "Validation failed", nil, errs)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	JSON(c, http.StatusUnauthorized, false, message, nil, nil)
}

func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal Server Error"
	}
	JSON(c, http.StatusInternalServerError, false, message, nil, nil)
}
