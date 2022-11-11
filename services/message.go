package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/logger"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/repositories"
	"github.com/rs/zerolog/log"
)

type IMessageService interface {
	FindByRoomID(roomID, userID string, params dto.MessageQueryParams) ([]models.Message, error)
	Create(roomID, userID string, dto dto.CreateMessage) (models.Message, error)
	Update(messageID, userID string, dto dto.UpdateMessage) (models.Message, error)
	Delete(messageID, userID string) (models.Message, error)
	CreateRemoveUser(roomID, removeUserID, userID string) (models.Message, error)
	CreateAddUser(roomID, addUserID, userID string) (models.Message, error)
	CreateCreateRoom(roomID, userID string) (models.Message, error)
	CreateUpdateRoom(roomID, userID string) (models.Message, error)
}

type MessageService struct {
	messageRepository repositories.IMessageRepository
	roomRepository    repositories.IRoomRepository
	userRepository    repositories.IUserRepository
}

var (
	messageOnce    sync.Once
	messageService IMessageService
)

func GetMessageService() IMessageService {
	messageOnce.Do(func() {
		log.Info().Msg("Initializing message service")
		messageService = &MessageService{
			messageRepository: repositories.GetMessageRepository(),
			roomRepository:    repositories.GetRoomRepository(),
			userRepository:    repositories.GetUserRepository(),
		}
	})
	return messageService
}

func (s *MessageService) FindByRoomID(roomID, userID string, params dto.MessageQueryParams) ([]models.Message, error) {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, userID).Msg("Finding messages")

	var messages []models.Message

	room, err := s.roomRepository.FindByID(roomID)
	if err != nil {
		return messages, err
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return messages, err
	}

	err = s.verifyIfCanRead(room, user)
	if err != nil {
		return messages, err
	}

	messages = s.messageRepository.FindByRoomID(roomID, params.Page, params.Size)

	return messages, nil
}

func (s *MessageService) Create(roomID, userID string, dto dto.CreateMessage) (models.Message, error) {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, userID).Msg("Creating message")

	var message models.Message

	room, err := s.roomRepository.FindByID(roomID)
	if err != nil {
		return message, err
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return message, err
	}

	err = s.verifyIfCanWrite(room, user)
	if err != nil {
		return message, err
	}

	message.Text = dto.Text
	message.RoomID = roomID
	message.UserID = userID
	message.Command = models.CreateMessage
	message.TargetID = roomID

	err = s.messageRepository.Create(&message)
	if err != nil {
		return message, err
	}

	message.User = user

	return message, nil
}

func (s *MessageService) Update(messageID, userID string, dto dto.UpdateMessage) (models.Message, error) {
	log.Debug().Str(logger.MessageID, messageID).Str(logger.UserID, userID).Msg("Updating message")

	message, err := s.messageRepository.FindByID(messageID)
	if err != nil {
		return message, err
	}

	if message.UserID != userID {
		return message, errors.New("you are not allowed to update this message")
	}

	message.Text = dto.Text
	message.Command = models.UpdateMessage
	message.TargetID = messageID

	err = s.messageRepository.Update(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) Delete(messageID, userID string) (models.Message, error) {
	log.Debug().Str(logger.MessageID, messageID).Str(logger.UserID, userID).Msg("Deleting message")

	message, err := s.messageRepository.FindByID(messageID)
	if err != nil {
		return message, err
	}

	if message.UserID != userID {
		return message, errors.New("you are not allowed to delete this message")
	}

	message.Text = ""
	message.Command = models.DeleteMessage
	message.TargetID = messageID

	err = s.messageRepository.Update(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) CreateRemoveUser(roomID, removeUserID, userID string) (models.Message, error) {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, removeUserID).Msg("Creating remove user message")

	var message models.Message

	removeUser, err := s.userRepository.FindByID(removeUserID)
	if err != nil {
		return message, err
	}

	message.Text = fmt.Sprintf("%s left the room", removeUser.FirstName)
	message.RoomID = roomID
	message.UserID = userID
	message.Command = models.RemoveUser
	message.TargetID = removeUserID

	err = s.messageRepository.Create(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) CreateAddUser(roomID, addUserID, userID string) (models.Message, error) {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, addUserID).Msg("Creating add user message")

	var message models.Message

	addUser, err := s.userRepository.FindByID(addUserID)
	if err != nil {
		return message, err
	}

	message.Text = fmt.Sprintf("%s joined the room", addUser.FirstName)
	message.RoomID = roomID
	message.UserID = userID
	message.Command = models.AddUser
	message.TargetID = addUserID

	err = s.messageRepository.Create(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) CreateCreateRoom(roomID, userID string) (models.Message, error) {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, userID).Msg("Creating create room message")

	var message models.Message

	message.Text = "Room created"
	message.RoomID = roomID
	message.UserID = userID
	message.Command = models.CreateRoom
	message.TargetID = roomID

	err := s.messageRepository.Create(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) CreateUpdateRoom(roomID, userID string) (models.Message, error) {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, userID).Msg("Creating update room message")

	var message models.Message

	message.Text = "Room updated"
	message.RoomID = roomID
	message.UserID = userID
	message.Command = models.UpdateRoom
	message.TargetID = roomID

	err := s.messageRepository.Create(&message)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (s *MessageService) verifyIfCanWrite(room models.Room, user models.User) error {
	log.Debug().Str(logger.RoomID, room.ID).Str(logger.UserID, user.ID).Msg("Verifying if user is authorized in room")

	switch room.RoomType {
	case models.PublicRoom:
		return nil
	case models.ReadOnlyRoom:
		if user.HasRole(models.AdminRole) {
			return nil
		}
	}

	return s.roomRepository.VerifyUserInRoom(room.ID, user.ID)
}

func (s *MessageService) verifyIfCanRead(room models.Room, user models.User) error {
	log.Debug().Str(logger.RoomID, room.ID).Str(logger.UserID, user.ID).Msg("Verifying if user is authorized in room")

	switch room.RoomType {
	case models.PublicRoom:
		return nil
	case models.ReadOnlyRoom:
		return nil
	}

	return s.roomRepository.VerifyUserInRoom(room.ID, user.ID)
}
