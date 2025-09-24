package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xhaiwa/user-service-golang/internal/handler"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", handler.HealthHandler)
	return r
}

func TestHealthCheckIntegration(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
