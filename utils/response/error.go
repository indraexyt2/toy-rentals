package response

import "github.com/gin-gonic/gin"

type responseError struct {
	Status  int         `json:"status_code,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

func ResponseError(c *gin.Context, statusCode int, message interface{}) {
	response := responseError{
		Status:  -1,
		Message: message,
	}

	c.JSON(statusCode, response)
}
