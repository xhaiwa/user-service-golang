package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xhaiwa/user-service-golang/internal/handler"
	"github.com/xhaiwa/user-service-golang/internal/repository"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"
)

func main() {
	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate crée la table si elle n'existe pas
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database migrated successfully!")

	r := gin.Default()

	r.GET("/health", handler.HealthCheck)
	r.POST("/signup", handler.SignupHandler(db))

	r.GET("/users/:id", handler.GetUserByIdHandler)

	r.Run(":8080")
}
