package repositories

import (
	"sync"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindAll() []models.User
	FindByID(id string) (models.User, error)
	FindByIdWithRooms(id string) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}

type UserRepository struct {
	DB *gorm.DB
}

var (
	userOnce       sync.Once
	userRepository IUserRepository
)

func GetUserRepository() IUserRepository {
	userOnce.Do(func() {
		log.Info().Msg("Initializing user repository")
		userRepository = &UserRepository{
			DB: models.GetDB(),
		}
	})
	return userRepository
}

func (r *UserRepository) FindAll() []models.User {
	var users []models.User
	r.DB.Find(&users)
	return users
}

func (r *UserRepository) FindByID(id string) (models.User, error) {
	var user models.User
	err := r.DB.First(&user, "id = ?", id).Error

	return user, err
}

func (r *UserRepository) FindByIdWithRooms(id string) (models.User, error) {
	var user models.User
	err := r.DB.Model(&models.User{}).Preload("Rooms").First(&user, "id = ?", id).Error

	return user, err
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.DB.First(&user, "email = ?", email).Error

	return user, err
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}
