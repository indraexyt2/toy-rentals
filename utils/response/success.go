package response

import "github.com/gin-gonic/gin"

type responseSuccess struct {
	Status   int         `json:"status_code,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Message  string      `json:"message,omitempty"`
}

func ResponseSuccess(c *gin.Context, statusCode int, data, metadata interface{}, message string) {
	var newResponse = &responseSuccess{
		Status:   1,
		Data:     data,
		Metadata: metadata,
		Message:  message,
	}

	c.JSON(statusCode, newResponse)
}
