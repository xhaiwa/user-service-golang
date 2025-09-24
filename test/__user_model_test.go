package test

import (
	"testing"
	"time"

	"github.com/xhaiwa/user-service-golang/internal/repository"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := repository.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	// Supprime la table si elle existe pour tests propres
	_ = db.Migrator().DropTable(&models.User{})
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate User table: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)

	user := models.User{
		Email:    "test@example.com",
		Password: "hashedpassword123",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Fatal("Expected user ID to be set")
	}

	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		t.Fatal("Expected timestamps to be set")
	}
}

func TestUniqueEmail(t *testing.T) {
	db := setupTestDB(t)

	user1 := models.User{Email: "dup@example.com", Password: "pw1"}
	user2 := models.User{Email: "dup@example.com", Password: "pw2"}

	if err := db.Create(&user1).Error; err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	err := db.Create(&user2).Error
	if err == nil {
		t.Fatal("Expected error on duplicate email, got nil")
	}
}

func TestSoftDelete(t *testing.T) {
	db := setupTestDB(t)

	user := models.User{Email: "delete@example.com", Password: "pw"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Soft delete
	if err := db.Delete(&user).Error; err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify it's not returned in normal queries
	var count int64
	db.Model(&models.User{}).Where("email = ?", "delete@example.com").Count(&count)
	if count != 0 {
		t.Fatal("Expected soft deleted user to be excluded from queries")
	}

	// Verify it exists in DB with DeletedAt
	var deletedUser models.User
	if err := db.Unscoped().First(&deletedUser, user.ID).Error; err != nil {
		t.Fatalf("Expected to find soft deleted user with Unscoped, got error: %v", err)
	}

	if deletedUser.DeletedAt.Time.IsZero() {
		t.Fatal("Expected DeletedAt to be set")
	}
}

func TestUpdatedAtChanges(t *testing.T) {
	db := setupTestDB(t)

	user := models.User{Email: "update@example.com", Password: "pw"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	createdAt := user.CreatedAt
	time.Sleep(time.Millisecond * 10) // assure timestamp diff

	user.Password = "newhashedpw"
	if err := db.Save(&user).Error; err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if !user.UpdatedAt.After(createdAt) {
		t.Fatal("Expected UpdatedAt to be updated after modification")
	}
}
