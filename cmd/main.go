package main

import (
	"github.com/gin-gonic/gin"

	"github.com/xhaiwa/user-service-golang/internal/handler"
)

func main() {
	r := gin.Default()

	r.GET("/health", handler.HealthHandler)

	r.Run(":8080")
}
