package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xhaiwa/user-service-golang/internal/handler"
	"github.com/xhaiwa/user-service-golang/internal/repository"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", handler.HealthCheck)
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

func TestDBConnection(t *testing.T) {
	db, err := repository.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get sql.DB from gorm.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("DB ping failed: %v", err)
	}
}
