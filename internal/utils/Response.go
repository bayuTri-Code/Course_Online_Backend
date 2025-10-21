package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status  bool      `json:"status"`           
	Message string      `json:"message,omitempty"`
}

type StandardResponse struct {
	Status  bool      `json:"status"`           
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`   
	Errors  interface{} `json:"errors,omitempty"` 
}

type ErrorResponse struct {
	Status  bool      `json:"status" example:"false"`
	Message string      `json:"message" example:"Server Internal error 500"`
	Errors  interface{} `json:"errors,omitempty"`
}


func JSONSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, StandardResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func JSONCreated(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, StandardResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func JSONError(c *gin.Context, message string, code int, errors interface{}) {
	c.JSON(code, ErrorResponse{
		Status:  false,
		Message: message,
		Errors: errors,
	})
}


func JSONNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	c.JSON(http.StatusNotFound, StandardResponse{
		Status:  false,
		Message: message,
	})
}

func JSONUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	c.JSON(http.StatusUnauthorized, StandardResponse{
		Status:  false,
		Message: message,
	})
}

