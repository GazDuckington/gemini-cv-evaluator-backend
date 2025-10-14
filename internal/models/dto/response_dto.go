package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`            // "success" or "error"
	Message string      `json:"message,omitempty"` // optional human-readable message
	Data    interface{} `json:"data,omitempty"`    // payload for success
	Error   interface{} `json:"error,omitempty"`   // error details for failure
}

type PaginatedData struct {
	Items    interface{} `json:"items"`               // usually a slice (e.g. []UserResponse)
	Total    int         `json:"total"`               // total items across all pages
	Page     int         `json:"page,omitempty"`      // current page (optional)
	PageSize int         `json:"page_size,omitempty"` // page size (optional)
}

// Success sends a standardized success response
func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// Error sends a standardized error response
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
		Error:   err,
	})
}

func SuccessPaginated(c *gin.Context, statusCode int, items interface{}, total, page, pageSize int, message string) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	p := PaginatedData{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    p,
	})
}
