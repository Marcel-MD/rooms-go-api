package services

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/repositories"
	"github.com/rs/zerolog/log"
)

type IMessageService interface {
	FindByRoomID(roomID, userID string, params dto.MessageQueryParams) ([]models.Message, error)
	Create(roomID, userID string, dto dto.CreateMessage) (models.Message, error)
	Update(id, userID string, dto dto.UpdateMessage) (models.Message, error)
	Delete(id, userID string) error
	VerifyUserInRoom(roomID, userID string) error
}

type MessageService struct {
	MessageRepository repositories.IMessageRepository
	RoomRepository    repositories.IRoomRepository
	UserRepository    repositories.IUserRepository
}

var (
	messageOnce    sync.Once
	messageService IMessageService
)

func GetMessageService() IMessageService {
	messageOnce.Do(func() {
		log.Info().Msg("Initializing message service")
		messageService = &MessageService{
			MessageRepository: repositories.GetMessageRepository(),
			RoomRepository:    repositories.GetRoomRepository(),
			UserRepository:    repositories.GetUserRepository(),
		}
	})
	return messageService
}

func (s *MessageService) FindByRoomID(roomID, userID string, params dto.MessageQueryParams) ([]models.Message, error) {
	log.Debug().Str("room_id", roomID).Str("user_id", userID).Msg("Finding messages")

	var messages []models.Message

	err := s.VerifyUserInRoom(roomID, userID)
	if err != nil {
		return messages, err
	}

	messages = s.MessageRepository.FindByRoomID(roomID, params.Page, params.Size)

	return messages, nil
}

func (s *MessageService) Create(roomID, userID string, dto dto.CreateMessage) (models.Message, error) {
	log.Debug().Str("room_id", roomID).Str("user_id", userID).Msg("Creating message")

	var message models.Message
	err := s.VerifyUserInRoom(roomID, userID)
	if err != nil {
		return message, err
	}

	user, err := s.UserRepository.FindByID(userID)
	if err != nil {
		return message, err
	}

	message.Text = dto.Text
	message.RoomID = roomID
	message.UserID = userID

	err = s.MessageRepository.Create(&message)
	if err != nil {
		return message, err
	}

	message.User = user

	return message, nil
}

func (s *MessageService) Update(id, userID string, dto dto.UpdateMessage) (models.Message, error) {
	log.Debug().Str("id", id).Str("user_id", userID).Msg("Updating message")

	message, err := s.MessageRepository.FindByID(id)
	if err != nil {
		return message, err
	}

	if message.UserID != userID {
		return message, errors.New("you are not allowed to update this message")
	}

	message.Text = dto.Text

	err = s.MessageRepository.Update(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) Delete(id, userID string) error {
	log.Debug().Str("id", id).Str("user_id", userID).Msg("Deleting message")

	message, err := s.MessageRepository.FindByID(id)
	if err != nil {
		return err
	}

	if message.UserID != userID {
		return errors.New("you are not allowed to delete this message")
	}

	err = s.MessageRepository.Delete(&message)
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) VerifyUserInRoom(roomID, userID string) error {
	log.Debug().Str("room_id", roomID).Str("user_id", userID).Msg("Verifying user in room")

	room, err := s.RoomRepository.FindByIdWithUsers(roomID)
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
