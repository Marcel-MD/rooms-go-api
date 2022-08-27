package services

import (
	"errors"
	"fmt"
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
	SearchByEmail(email string) []models.User
	FindOne(id string) (models.User, error)
	SendOtp(email string) error
	RegisterOtp(dto dto.RegisterOtpUser) (models.User, error)
	Register(dto dto.RegisterUser) (models.User, error)
	LoginOtp(dto dto.LoginOtpUser) (string, error)
	Login(dto dto.LoginUser) (string, error)
	Update(dto dto.UpdateUser, id string) (models.User, error)
}

type UserService struct {
	repository  repositories.IUserRepository
	otpService  IOtpService
	mailService IMailService
}

var (
	userOnce    sync.Once
	userService IUserService
)

func GetUserService() IUserService {
	userOnce.Do(func() {
		log.Info().Msg("Initializing user service")
		userService = &UserService{
			repository:  repositories.GetUserRepository(),
			otpService:  GetOtpService(),
			mailService: GetMailService(),
		}
	})
	return userService
}

func (s *UserService) FindAll() []models.User {
	log.Debug().Msg("Finding all users")

	return s.repository.FindAll()
}

func (s *UserService) SearchByEmail(email string) []models.User {
	log.Debug().Msg("Searching for users by email")

	return s.repository.SearchByEmail(email)
}

func (s *UserService) FindOne(id string) (models.User, error) {
	log.Debug().Str("id", id).Msg("Finding user")

	user, err := s.repository.FindByIdWithRooms(id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) SendOtp(email string) error {
	log.Debug().Msg("Sending otp")

	otp, err := s.otpService.Generate(email)
	if err != nil {
		return err
	}

	mail := Mail{
		To:      []string{email},
		Subject: "Rooms - Verification Code",
		Body:    fmt.Sprintf("Your verification code is <strong>%s</strong>.", otp),
	}

	go s.mailService.Send(mail)

	return nil
}

func (s *UserService) RegisterOtp(dto dto.RegisterOtpUser) (models.User, error) {
	log.Debug().Msg("Registering user with otp")

	var user models.User

	err := s.otpService.Verify(dto.Email, dto.Otp)
	if err != nil {
		return user, err
	}

	return s.Register(dto.RegisterUser)
}

func (s *UserService) Register(dto dto.RegisterUser) (models.User, error) {
	log.Debug().Msg("Registering user")

	user, err := s.repository.FindByEmail(dto.Email)
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

	err = s.repository.Create(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) LoginOtp(dto dto.LoginOtpUser) (string, error) {
	log.Debug().Msg("Logging in user with otp")

	err := s.otpService.Verify(dto.Email, dto.Otp)
	if err != nil {
		return "", err
	}

	return s.Login(dto.LoginUser)
}

func (s *UserService) Login(dto dto.LoginUser) (string, error) {
	log.Debug().Msg("Logging in user")

	user, err := s.repository.FindByEmail(dto.Email)
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

	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName

	err = s.repository.Update(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
