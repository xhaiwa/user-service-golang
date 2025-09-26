package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xhaiwa/user-service-golang/internal/handler"
	"github.com/xhaiwa/user-service-golang/internal/repository"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var integrationDB *gorm.DB

func TestMain(m *testing.M) {
	db, terminate := setupPostgresContainer()
	integrationDB = db
	code := m.Run()
	terminate()
	os.Exit(code)
}

func setupPostgresContainer() (*gorm.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_SECRET", "testsecret")

	db, err := repository.ConnectDB()
	if err != nil {
		panic(err)
	}

	return db, func() {
		container.Terminate(ctx)
	}
}

func resetTables(db *gorm.DB) {
	db.Exec("DROP TABLE IF EXISTS users")
	db.AutoMigrate(&models.User{})
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	return r
}

func TestSignupIntegration(t *testing.T) {
	resetTables(integrationDB)
	router := setupRouter()
	router.POST("/signup", handler.SignupHandler(integrationDB))

	t.Run("Create user successfully", func(t *testing.T) {
		payload := map[string]string{
			"email":    "integration@example.com",
			"password": "strongpassword",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 201, resp.Code)
		assert.Contains(t, resp.Body.String(), "token")
	})

	t.Run("Duplicate email", func(t *testing.T) {
		payload := map[string]string{
			"email":    "integration@example.com",
			"password": "strongpassword",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 400, resp.Code)
		assert.Contains(t, resp.Body.String(), "email already exists")
	})
}
