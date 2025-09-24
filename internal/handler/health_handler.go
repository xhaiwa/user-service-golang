package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xhaiwa/user-service-golang/internal/repository"
)

func HealthCheck(c *gin.Context) {
	_, err := repository.ConnectDB()
	dbStatus := "ok"
	if err != nil {
		dbStatus = "error: " + err.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "User Service is running",
		"db":     dbStatus,
	})
}
