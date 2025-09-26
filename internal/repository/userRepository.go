package repository

import (
	"errors"

	"github.com/xhaiwa/user-service-golang/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, email, password string) (*models.User, error) {
	// Vérifier si email existe
	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	var oauthID *string = nil
	user := models.User{
		Email:         email,
		Password:      hashed,
		OAuthID:       oauthID,
		OAuthProvider: nil,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// HashPassword hash le mot de passe
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
