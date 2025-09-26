package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xhaiwa/user-service-golang/internal/repository/models"
	"github.com/xhaiwa/user-service-golang/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword hash le mot de passe
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword compare le hash et le mot de passe
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateJWT génère un token JWT
func CreateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.JWTSecret))
}

// Signup crée un utilisateur dans la DB
func Signup(db *gorm.DB, email, password string) (*models.User, string, error) {
	// Vérifier si email existe
	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, "", errors.New("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, "", err
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	var oauthID *string = nil
	user := models.User{
		Email:         email,
		Password:      hashed,
		OAuthID:       oauthID,
		OAuthProvider: nil,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, "", err
	}

	// Générer JWT
	token, err := CreateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}
