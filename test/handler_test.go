package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/xhaiwa/user-service-golang/internal/handler"
	"github.com/xhaiwa/user-service-golang/internal/repository"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"
)

func setupHandlerDB(t *testing.T) *gorm.DB {
	db, err := repository.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect DB: %v", err)
	}
	db.Exec("DROP TABLE IF EXISTS users")
	db.AutoMigrate(&models.User{})
	return db
}

func TestSignupHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupHandlerDB(t)

	router := gin.Default()
	router.POST("/signup", handler.SignupHandler(db))

	payload := map[string]string{"email": "test@example.com", "password": "pwd123"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
	assert.Contains(t, resp.Body.String(), "test@example.com")
}

func TestHealthCheckHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/health", handler.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "ok")
}
