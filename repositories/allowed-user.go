package repositories

import (
	"sync"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type IAllowedUserRepository interface {
	FindAll() []models.AllowedUser
	FindByEmail(email string) (models.AllowedUser, error)
	Create(allowedUser *models.AllowedUser) error
	Delete(email string) error
}

type AllowedUserRepository struct {
	DB *gorm.DB
}

var (
	allowedUserOnce       sync.Once
	allowedUserRepository IAllowedUserRepository
)

func GetAllowedUserRepository() IAllowedUserRepository {
	allowedUserOnce.Do(func() {
		log.Info().Msg("Initializing allowed user repository")
		allowedUserRepository = &AllowedUserRepository{
			DB: models.GetDB(),
		}
	})
	return allowedUserRepository
}

func (r *AllowedUserRepository) FindAll() []models.AllowedUser {
	var allowedUsers []models.AllowedUser
	r.DB.Find(&allowedUsers)
	return allowedUsers
}

func (r *AllowedUserRepository) FindByEmail(email string) (models.AllowedUser, error) {
	var allowedUser models.AllowedUser
	err := r.DB.First(&allowedUser, "email = ?", email).Error

	return allowedUser, err
}

func (r *AllowedUserRepository) Create(allowedUser *models.AllowedUser) error {
	return r.DB.Create(allowedUser).Error
}

func (r *AllowedUserRepository) Delete(email string) error {
	return r.DB.Delete(&models.AllowedUser{}, "email = ?", email).Error
}
