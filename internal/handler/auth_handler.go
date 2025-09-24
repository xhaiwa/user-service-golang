package handler

import (
	"net/http"

	"github.com/xhaiwa/user-service-golang/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func SignupHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, token, err := service.Signup(db, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
			"token": token,
		})
	}
}
