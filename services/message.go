package services

import (
	"errors"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"gorm.io/gorm"
)

type IMessageService interface {
	FindByRoomID(roomID string, userID string, params dto.MessageQueryParams) ([]models.Message, error)
	Create(dto dto.CreateMessage, roomID string, userID string) (models.Message, error)
	Update(id string, dto dto.UpdateMessage, userID string) (models.Message, error)
	Delete(id string, userID string) error
}

type MessageService struct {
	DB *gorm.DB
}

func NewMessageService() IMessageService {
	return &MessageService{
		DB: models.GetDB(),
	}
}

func (s *MessageService) FindByRoomID(roomID string, userID string, params dto.MessageQueryParams) ([]models.Message, error) {

	var messages []models.Message

	err := s.verifyUserInRoom(roomID, userID)
	if err != nil {
		return messages, err
	}

	s.DB.Scopes(models.Paginate(params.Page, params.Size)).Model(&models.Message{}).Order("created_at desc").Preload("User").Find(&messages, "room_id = ?", roomID)
	return messages, nil
}

func (s *MessageService) Create(dto dto.CreateMessage, roomID string, userID string) (models.Message, error) {

	var message models.Message

	err := s.verifyUserInRoom(roomID, userID)
	if err != nil {
		return message, err
	}

	message.Text = dto.Text
	message.RoomID = roomID
	message.UserID = userID

	err = s.DB.Create(&message).Error
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) Update(id string, dto dto.UpdateMessage, userID string) (models.Message, error) {

	var message models.Message

	err := s.DB.First(&message, "id = ?", id).Error
	if err != nil {
		return message, err
	}

	if message.UserID != userID {
		return message, errors.New("you are not allowed to update this message")
	}

	message.Text = dto.Text

	err = s.DB.Save(&message).Error
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) Delete(id string, userID string) error {

	var message models.Message

	err := s.DB.First(&message, "id = ?", id).Error
	if err != nil {
		return err
	}

	if message.UserID != userID {
		return errors.New("you are not allowed to delete this message")
	}

	err = s.DB.Delete(&message).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) verifyUserInRoom(roomID string, userID string) error {
	var room models.Room
	err := s.DB.Model(&models.Room{}).Preload("Users").First(&room, "id = ?", roomID).Error
	if err != nil {
		return err
	}

	for _, user := range room.Users {
		if user.ID == userID {
			return nil
		}
	}

	return errors.New("user is not in room")
}
