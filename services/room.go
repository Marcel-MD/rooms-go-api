package services

import (
	"errors"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"gorm.io/gorm"
)

type IRoomService interface {
	FindAll() []models.Room
	FindOne(id string) (models.Room, error)
	Create(dto dto.CreateRoom, userID string) (models.Room, error)
	Update(id string, dto dto.UpdateRoom, userID string) (models.Room, error)
	Delete(id string, userID string) error
	AddUser(id string, email string, userID string) error
	RemoveUser(id string, email string, userID string) error
}

type RoomService struct {
	DB *gorm.DB
}

func NewRoomService() IRoomService {
	return &RoomService{
		DB: models.GetDB(),
	}
}

func (s *RoomService) FindAll() []models.Room {
	var rooms []models.Room
	s.DB.Find(&rooms)
	return rooms
}

func (s *RoomService) FindOne(id string) (models.Room, error) {
	var room models.Room

	err := s.DB.Model(&models.Room{}).Preload("Users").First(&room, "id = ?", id).Error
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Create(dto dto.CreateRoom, userID string) (models.Room, error) {
	var user models.User
	err := s.DB.First(&user, "id = ?", userID).Error

	var room models.Room
	room = models.Room{
		Name:    dto.Name,
		OwnerID: userID,
	}

	err = s.DB.Create(&room).Error
	if err != nil {
		return room, err
	}

	err = s.AddUser(room.ID, user.Email, userID)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Update(id string, dto dto.UpdateRoom, userID string) (models.Room, error) {
	var room models.Room

	err := s.DB.First(&room, "id = ?", id).Error
	if err != nil {
		return room, err
	}

	if room.OwnerID != userID {
		return room, errors.New("you are not the owner of this room")
	}

	room.Name = dto.Name

	err = s.DB.Save(&room).Error
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Delete(id string, userID string) error {
	var room models.Room

	err := s.DB.First(&room, "id = ?", id).Error
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return errors.New("you are not the owner of this room")
	}

	err = s.DB.Delete(&room).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) AddUser(id string, email string, userID string) error {
	var room models.Room

	err := s.DB.First(&room, "id = ?", id).Error
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return errors.New("you are not the owner of this room")
	}

	var user models.User

	err = s.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return err
	}

	err = s.DB.Model(&room).Omit("Users.*").Association("Users").Append(&user)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) RemoveUser(id string, email string, userID string) error {
	var room models.Room

	err := s.DB.First(&room, "id = ?", id).Error
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return errors.New("you are not the owner of this room")
	}

	var user models.User

	err = s.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return err
	}

	err = s.DB.Model(&room).Association("Users").Delete(&user)
	if err != nil {
		return err
	}

	return nil
}
