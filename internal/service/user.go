package service

import (
	"errors"

	"github.com/xhaiwa/user-service-golang/internal/repository/models"
	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, id int) (*models.User, error) {
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // ou custom NotFoundError
		}
		return nil, err
	}
	return &user, nil
}
