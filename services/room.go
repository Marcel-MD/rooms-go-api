package services

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/repositories"
	"github.com/rs/zerolog/log"
)

type IRoomService interface {
	FindAll() []models.Room
	FindOne(id string) (models.Room, error)
	Create(dto dto.CreateRoom, userID string) (models.Room, error)
	Update(id, userID string, dto dto.UpdateRoom) (models.Room, error)
	Delete(id, userID string) error
	AddUser(id, email, userID string) error
	RemoveUser(id, removeUserID, userID string) error
}

type RoomService struct {
	RoomRepository repositories.IRoomRepository
	UserRepository repositories.IUserRepository
}

var (
	roomOnce    sync.Once
	roomService IRoomService
)

func GetRoomService() IRoomService {
	roomOnce.Do(func() {
		log.Info().Msg("Initializing room service")
		roomService = &RoomService{
			RoomRepository: repositories.GetRoomRepository(),
			UserRepository: repositories.GetUserRepository(),
		}
	})
	return roomService
}

func (s *RoomService) FindAll() []models.Room {
	log.Debug().Msg("Finding all rooms")

	return s.RoomRepository.FindAll()
}

func (s *RoomService) FindOne(id string) (models.Room, error) {
	log.Debug().Str("id", id).Msg("Finding room")

	room, err := s.RoomRepository.FindByIdWithUsers(id)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Create(dto dto.CreateRoom, userID string) (models.Room, error) {
	log.Debug().Str("user_id", userID).Msg("Creating room")

	user, err := s.UserRepository.FindByID(userID)
	if err != nil {
		return models.Room{}, err
	}

	room := models.Room{
		Name:    dto.Name,
		OwnerID: userID,
	}

	err = s.RoomRepository.Create(&room)
	if err != nil {
		return room, err
	}

	err = s.AddUser(room.ID, user.Email, userID)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Update(id, userID string, dto dto.UpdateRoom) (models.Room, error) {
	log.Debug().Str("id", id).Str("user_id", userID).Msg("Updating room")

	room, err := s.RoomRepository.FindByID(id)
	if err != nil {
		return room, err
	}

	if room.OwnerID != userID {
		return room, errors.New("you are not the owner of this room")
	}

	room.Name = dto.Name

	err = s.RoomRepository.Update(&room)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Delete(id, userID string) error {
	log.Debug().Str("id", id).Str("user_id", userID).Msg("Deleting room")

	room, err := s.RoomRepository.FindByID(id)
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return errors.New("you are not the owner of this room")
	}

	err = s.RoomRepository.Delete(&room)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) AddUser(id, email, userID string) error {
	log.Debug().Str("id", id).Msg("Adding user to room")

	room, err := s.RoomRepository.FindByID(id)
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return errors.New("you are not the owner of this room")
	}

	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return err
	}

	err = s.RoomRepository.AddUser(&room, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) RemoveUser(roomId, removeUserID, userID string) error {
	log.Debug().Str("room_id", roomId).Str("user_id", removeUserID).Msg("Removing user from room")

	room, err := s.RoomRepository.FindByID(roomId)
	if err != nil {
		return err
	}

	user, err := s.UserRepository.FindByID(removeUserID)
	if err != nil {
		return err
	}

	if room.OwnerID != userID && removeUserID != userID {
		return errors.New("unauthorized")
	}

	if room.OwnerID == removeUserID {
		return errors.New("you are the owner of this room")
	}

	err = s.RoomRepository.RemoveUser(&room, &user)
	if err != nil {
		return err
	}

	return nil
}
