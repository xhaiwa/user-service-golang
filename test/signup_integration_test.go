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

	"gorm.io/gorm"

	"github.com/xhaiwa/user-service-golang/internal/handler"
	"github.com/xhaiwa/user-service-golang/internal/repository"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainer(t *testing.T) (*gorm.DB, func()) {
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
	assert.NoError(t, err)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("JWT_SECRET", "testsecret")

	db, err := repository.ConnectDB()
	assert.NoError(t, err)

	return db, func() {
		container.Terminate(ctx)
	}
}

func TestSignupIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, terminate := setupPostgresContainer(t)
	defer terminate()

	// Clean table
	db.Exec("DROP TABLE IF EXISTS users")
	db.AutoMigrate(&models.User{})

	router := gin.Default()
	router.POST("/signup", handler.SignupHandler(db))

	// 1️⃣ Création d'un utilisateur
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
	assert.Contains(t, resp.Body.String(), "integration@example.com")

	// 2️⃣ Email dupliqué
	req2, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, 400, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "email already exists")

	// 3️⃣ Vérifier que le mot de passe est hashé
	var user models.User
	db.Where("email = ?", "integration@example.com").First(&user)
	assert.NotEqual(t, "strongpassword", user.Password)
}
