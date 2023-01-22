package services

import (
	"errors"
	"fmt"
	"os"
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
	UpdateOtp(dto dto.UpdateOtpUser, id string) (models.User, error)
	Update(dto dto.UpdateUser, id string) (models.User, error)
	AddRole(id string, role string, userID string) (models.User, error)
	RemoveRole(id string, role string, userID string) (models.User, error)
	SetIsOnline(id string, isOnline bool) (models.User, error)
}

type UserService struct {
	repository            repositories.IUserRepository
	allowedUserRepository repositories.IAllowedUserRepository
	otpService            IOtpService
	mailService           IMailService
	loginLimiterService   ILoginLimiterService
	allowedUsersOnly      bool
}

var (
	userOnce    sync.Once
	userService IUserService
)

func GetUserService() IUserService {
	userOnce.Do(func() {
		log.Info().Msg("Initializing user service")
		us := &UserService{
			repository:            repositories.GetUserRepository(),
			allowedUserRepository: repositories.GetAllowedUserRepository(),
			otpService:            GetOtpService(),
			mailService:           GetMailService(),
			loginLimiterService:   GetLoginLimiterService(),
		}

		allowedUsersOnly := os.Getenv("ALLOWED_USERS_ONLY")
		if allowedUsersOnly == "true" {
			us.allowedUsersOnly = true
		}

		userService = us
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

	if s.allowedUsersOnly {
		_, err := s.allowedUserRepository.FindByEmail(email)
		if err != nil {
			return errors.New("user not allowed")
		}
	}

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
		Phone:     dto.Phone,
		Password:  string(hashedPassword),
		Roles:     []string{models.UserRole},
	}

	if s.allowedUsersOnly {
		allowedUser, err := s.allowedUserRepository.FindByEmail(dto.Email)
		if err != nil {
			return user, errors.New("user not allowed")
		}

		user.Roles = append(user.Roles, allowedUser.DefaultRole)
	}

	err = s.repository.Create(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) LoginOtp(dto dto.LoginOtpUser) (string, error) {
	log.Debug().Msg("Logging in user with otp")

	err := s.loginLimiterService.IncrementAttempts(dto.Email)
	if err != nil {
		return "", err
	}

	err = s.otpService.Verify(dto.Email, dto.Otp)
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

func (s *UserService) UpdateOtp(dto dto.UpdateOtpUser, id string) (models.User, error) {
	log.Debug().Msg("Updating user with otp")

	err := s.otpService.Verify(dto.Email, dto.Otp)
	if err != nil {
		return models.User{}, err
	}

	return s.Update(dto.UpdateUser, id)
}

func (s *UserService) Update(dto dto.UpdateUser, id string) (models.User, error) {
	log.Debug().Msg("Updating user")

	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.Email = dto.Email
	user.Phone = dto.Phone

	err = s.repository.Update(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) AddRole(id string, role string, userID string) (models.User, error) {
	log.Debug().Msg("Adding role to user")

	admin, err := s.repository.FindByID(userID)
	if err != nil {
		return admin, err
	}

	if !admin.HasRole(models.AdminRole) {
		return admin, errors.New("user is not admin")
	}

	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	if user.HasRole(role) {
		return user, errors.New("user already has this role")
	}

	user.Roles = append(user.Roles, role)

	err = s.repository.Update(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) RemoveRole(id string, role string, userID string) (models.User, error) {
	log.Debug().Msg("Removing role from user")

	admin, err := s.repository.FindByID(userID)
	if err != nil {
		return admin, err
	}

	if !admin.HasRole(models.AdminRole) {
		return admin, errors.New("user is not admin")
	}

	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	if !user.HasRole(role) {
		return user, errors.New("user does not have this role")
	}

	user.Roles = remove(user.Roles, role)

	err = s.repository.Update(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserService) SetIsOnline(id string, isOnline bool) (models.User, error) {

	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	user.IsOnline = isOnline

	err = s.repository.Update(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (s *UserService) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
