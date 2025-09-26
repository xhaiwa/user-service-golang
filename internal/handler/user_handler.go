package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserByIdHandler(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"user": id,
	})
}
