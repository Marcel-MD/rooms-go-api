package services

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/repositories"
	"github.com/Marcel-MD/rooms-go-api/token"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	FindAll() []models.User
	FindOne(id string) (models.User, error)
	Register(dto dto.RegisterUser) (models.User, error)
	Login(dto dto.LoginUser) (string, error)
	Update(dto dto.UpdateUser, id string) (models.User, error)
}

type UserService struct {
	Repository repositories.IUserRepository
}

var (
	userOnce    sync.Once
	userService IUserService
)

func GetUserService() IUserService {
	userOnce.Do(func() {
		log.Info().Msg("Initializing user service")
		userService = &UserService{
			Repository: repositories.GetUserRepository(),
		}
	})
	return userService
}

func (s *UserService) FindAll() []models.User {
	log.Debug().Msg("Finding all users")

	return s.Repository.FindAll()
}

func (s *UserService) FindOne(id string) (models.User, error) {
	log.Debug().Str("id", id).Msg("Finding user")

	user, err := s.Repository.FindByIdWithRooms(id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) Register(dto dto.RegisterUser) (models.User, error) {
	log.Debug().Msg("Registering user")

	user, err := s.Repository.FindByEmail(dto.Email)
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

	err = s.Repository.Create(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) Login(dto dto.LoginUser) (string, error) {
	log.Debug().Msg("Logging in user")

	user, err := s.Repository.FindByEmail(dto.Email)
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

	user, err := s.Repository.FindByID(id)
	if err != nil {
		return user, err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName

	err = s.Repository.Update(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
