package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xhaiwa/user-service-golang/internal/handler"
	"github.com/xhaiwa/user-service-golang/internal/repository"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"

	"github.com/xhaiwa/user-service-golang/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignupHandler(t *testing.T) {
	utils.JWTSecret = "testsecret" // charger un secret pour tests

	db, err := repository.ConnectDB()
	assert.NoError(t, err)

	db.AutoMigrate(&models.User{})

	// nettoyer DB avant test
	db.Exec("DELETE FROM users")

	router := gin.Default()
	router.POST("/signup", handler.SignupHandler(db))

	// 1️⃣ Test création d'un nouvel utilisateur
	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
	assert.Contains(t, resp.Body.String(), "test@example.com")

	// 2️⃣ Test email dupliqué
	req2, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, 400, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "email already exists")

	// 3️⃣ Test hash password
	var user models.User
	db.Where("email = ?", "test@example.com").First(&user)
	assert.NotEqual(t, "password123", user.Password)
}
