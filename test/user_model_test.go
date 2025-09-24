package test

import (
	"context"
	"testing"
	"time"

	"github.com/xhaiwa/user-service-golang/internal/repository"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
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

	// Configurer l'environnement pour ConnectDB
	t.Setenv("DB_HOST", host)
	t.Setenv("DB_PORT", port.Port())
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "password")
	t.Setenv("DB_NAME", "test_db")

	db, err := repository.ConnectDB()
	assert.NoError(t, err)

	// Supprimer la table pour test propre
	_ = db.Migrator().DropTable(&models.User{})
	assert.NoError(t, db.AutoMigrate(&models.User{}))

	return db, func() {
		container.Terminate(ctx)
	}
}

func TestCreateUser(t *testing.T) {
	db, terminate := setupTestDB(t)
	defer terminate()

	user := models.User{
		Email:    "test@example.com",
		Password: "hashedpassword123",
	}

	assert.NoError(t, db.Create(&user).Error)
	assert.NotZero(t, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUniqueEmail(t *testing.T) {
	db, terminate := setupTestDB(t)
	defer terminate()

	user1 := models.User{Email: "dup@example.com", Password: "pw1"}
	user2 := models.User{Email: "dup@example.com", Password: "pw2"}

	assert.NoError(t, db.Create(&user1).Error)
	assert.Error(t, db.Create(&user2).Error)
}

func TestSoftDelete(t *testing.T) {
	db, terminate := setupTestDB(t)
	defer terminate()

	user := models.User{Email: "delete@example.com", Password: "pw"}
	assert.NoError(t, db.Create(&user).Error)

	assert.NoError(t, db.Delete(&user).Error)

	var count int64
	db.Model(&models.User{}).Where("email = ?", "delete@example.com").Count(&count)
	assert.Equal(t, int64(0), count)

	var deletedUser models.User
	assert.NoError(t, db.Unscoped().First(&deletedUser, user.ID).Error)
	assert.False(t, deletedUser.DeletedAt.Time.IsZero())
}

func TestUpdatedAtChanges(t *testing.T) {
	db, terminate := setupTestDB(t)
	defer terminate()

	user := models.User{Email: "update@example.com", Password: "pw"}
	assert.NoError(t, db.Create(&user).Error)

	createdAt := user.CreatedAt
	time.Sleep(time.Millisecond * 10)

	user.Password = "newhashedpw"
	assert.NoError(t, db.Save(&user).Error)
	assert.True(t, user.UpdatedAt.After(createdAt))
}
