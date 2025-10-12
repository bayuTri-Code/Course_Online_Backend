package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Status  string      `json:"status"`           
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`   
	Errors  interface{} `json:"errors,omitempty"` 
}


func JSONSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, StandardResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func JSONCreated(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, StandardResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func JSONError(c *gin.Context, message string, code int) {
	c.JSON(code, StandardResponse{
		Status:  "error",
		Message: message,
	})
}


func JSONNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	c.JSON(http.StatusNotFound, StandardResponse{
		Status:  "error",
		Message: message,
	})
}

func JSONUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	c.JSON(http.StatusUnauthorized, StandardResponse{
		Status:  "error",
		Message: message,
	})
}

