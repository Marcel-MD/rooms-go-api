package services

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/token"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	FindAll() []models.User
	FindOne(id string) (models.User, error)
	Register(dto dto.RegisterUser) (models.User, error)
	Login(dto dto.LoginUser) (string, error)
	Update(dto dto.UpdateUser, id string) (models.User, error)
}

type UserService struct {
	DB *gorm.DB
}

var (
	userOnce    sync.Once
	userService IUserService
)

func GetUserService() IUserService {
	userOnce.Do(func() {
		log.Info().Msg("Initializing user service")
		userService = &UserService{
			DB: models.GetDB(),
		}
	})
	return userService
}

func (s *UserService) FindAll() []models.User {
	log.Debug().Msg("Finding all users")

	var users []models.User
	s.DB.Find(&users)
	return users
}

func (s *UserService) FindOne(id string) (models.User, error) {
	log.Debug().Str("id", id).Msg("Finding user")

	var user models.User
	err := s.DB.Model(&models.User{}).Preload("Rooms").First(&user, "id = ?", id).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) Register(dto dto.RegisterUser) (models.User, error) {
	log.Debug().Msg("Registering user")

	var user models.User
	err := s.DB.First(&user, "email = ?", dto.Email).Error
	if err == nil {
		return user, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user = models.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  string(hashedPassword),
	}

	err = s.DB.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) Login(dto dto.LoginUser) (string, error) {
	log.Debug().Msg("Logging in user")

	var user models.User
	err := s.DB.First(&user, "email = ?", dto.Email).Error
	if err != nil {
		return "", err
	}

	err = s.verifyPassword(dto.Password, user.Password)
	if err != nil {
		return "", err
	}

	return token.Generate(user.ID)
}

func (s *UserService) Update(dto dto.UpdateUser, id string) (models.User, error) {
	log.Debug().Msg("Updating user")

	var user models.User
	err := s.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return user, err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName

	err = s.DB.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
