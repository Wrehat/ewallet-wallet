package response

import "github.com/gin-gonic/gin"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
