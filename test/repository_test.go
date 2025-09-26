package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"
)

func TestRepository_CreateUser(t *testing.T) {
	resetTables(integrationDB)

	user := models.User{Email: "repo@test.com", Password: "pw"}
	err := integrationDB.Create(&user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestRepository_UniqueEmail(t *testing.T) {
	resetTables(integrationDB)

	user1 := models.User{Email: "dup@test.com", Password: "pw1"}
	user2 := models.User{Email: "dup@test.com", Password: "pw2"}

	assert.NoError(t, integrationDB.Create(&user1).Error)
	err := integrationDB.Create(&user2).Error
	assert.NotNil(t, err)
}

func TestRepository_SoftDelete(t *testing.T) {
	resetTables(integrationDB)

	user := models.User{Email: "delete@test.com", Password: "pw"}
	integrationDB.Create(&user)

	integrationDB.Delete(&user)

	var count int64
	integrationDB.Model(&models.User{}).Where("email = ?", "delete@test.com").Count(&count)
	assert.Equal(t, int64(0), count)

	var deletedUser models.User
	assert.NoError(t, integrationDB.Unscoped().First(&deletedUser, user.ID).Error)
	assert.False(t, deletedUser.DeletedAt.Time.IsZero())
}

func TestRepository_UpdatedAt(t *testing.T) {
	resetTables(integrationDB)

	user := models.User{Email: "update@test.com", Password: "pw"}
	integrationDB.Create(&user)

	createdAt := user.CreatedAt
	time.Sleep(10 * time.Millisecond)

	user.Password = "newpw"
	integrationDB.Save(&user)

	assert.True(t, user.UpdatedAt.After(createdAt))
}
