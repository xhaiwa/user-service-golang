package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xhaiwa/user-service-golang/internal/handler"

	"github.com/gin-gonic/gin"
)

func TestHealthCheckHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/health", handler.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	expectedBody := `"status":"User Service is running"`
	if !contains(w.Body.String(), expectedBody) {
		t.Errorf("Expected body to contain %s, got %s", expectedBody, w.Body.String())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
