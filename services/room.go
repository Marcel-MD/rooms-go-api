package services

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/dto"
	"github.com/Marcel-MD/rooms-go-api/logger"
	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/Marcel-MD/rooms-go-api/repositories"
	"github.com/rs/zerolog/log"
)

type IRoomService interface {
	FindAll() []models.Room
	FindOne(id string) (models.Room, error)
	Create(dto dto.CreateRoom, userID string) (models.Room, error)
	Update(roomID, userID string, dto dto.UpdateRoom) (models.Room, error)
	Delete(roomID, userID string) error
	AddUser(roomID, addUserID, userID string) error
	RemoveUser(roomID, removeUserID, userID string) error
	VerifyUserInRoom(roomID, userID string) error
	createDefaultRooms()
}

type RoomService struct {
	roomRepository repositories.IRoomRepository
	userRepository repositories.IUserRepository
}

var (
	roomOnce    sync.Once
	roomService IRoomService
)

func GetRoomService() IRoomService {
	roomOnce.Do(func() {
		log.Info().Msg("Initializing room service")
		roomService = &RoomService{
			roomRepository: repositories.GetRoomRepository(),
			userRepository: repositories.GetUserRepository(),
		}
		roomService.createDefaultRooms()
	})
	return roomService
}

func (s *RoomService) FindAll() []models.Room {
	log.Debug().Msg("Finding all rooms")

	return s.roomRepository.FindAll()
}

func (s *RoomService) FindOne(id string) (models.Room, error) {
	log.Debug().Str(logger.RoomID, id).Msg("Finding room")

	room, err := s.roomRepository.FindByIdWithUsers(id)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Create(dto dto.CreateRoom, userID string) (models.Room, error) {
	log.Debug().Str(logger.UserID, userID).Msg("Creating room")

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return models.Room{}, err
	}

	room := models.Room{
		Name:     dto.Name,
		OwnerID:  userID,
		RoomType: models.PrivateRoom,
	}

	err = s.roomRepository.Create(&room)
	if err != nil {
		return room, err
	}

	err = s.AddUser(room.ID, user.ID, userID)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Update(roomID, userID string, dto dto.UpdateRoom) (models.Room, error) {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, userID).Msg("Updating room")

	room, err := s.roomRepository.FindByID(roomID)
	if err != nil {
		return room, err
	}

	if room.OwnerID != userID {
		return room, errors.New("you are not the owner of this room")
	}

	room.Name = dto.Name

	err = s.roomRepository.Update(&room)
	if err != nil {
		return room, err
	}

	return room, nil
}

func (s *RoomService) Delete(roomID, userID string) error {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, userID).Msg("Deleting room")

	room, err := s.roomRepository.FindByID(roomID)
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return errors.New("you are not the owner of this room")
	}

	err = s.roomRepository.Delete(&room)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) AddUser(roomID, addUserID, userID string) error {
	log.Debug().Str(logger.RoomID, roomID).Msg("Adding user to room")

	err := s.roomRepository.VerifyUserInRoom(roomID, addUserID)
	if err == nil {
		return errors.New("user already in this room")
	}

	room, err := s.roomRepository.FindByID(roomID)
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return errors.New("you are not the owner of this room")
	}

	user, err := s.userRepository.FindByID(addUserID)
	if err != nil {
		return err
	}

	err = s.roomRepository.AddUser(&room, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) RemoveUser(roomID, removeUserID, userID string) error {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, removeUserID).Msg("Removing user from room")

	err := s.roomRepository.VerifyUserInRoom(roomID, removeUserID)
	if err != nil {
		return err
	}

	room, err := s.roomRepository.FindByID(roomID)
	if err != nil {
		return err
	}

	user, err := s.userRepository.FindByID(removeUserID)
	if err != nil {
		return err
	}

	if room.OwnerID != userID && removeUserID != userID {
		return errors.New("unauthorized")
	}

	if room.OwnerID == removeUserID {
		return errors.New("you are the owner of this room")
	}

	err = s.roomRepository.RemoveUser(&room, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomService) VerifyUserInRoom(roomID, userID string) error {
	log.Debug().Str(logger.RoomID, roomID).Str(logger.UserID, userID).Msg("Verifying user in room")
	return s.roomRepository.VerifyUserInRoom(roomID, userID)
}

func (s *RoomService) createDefaultRooms() {
	_, err := s.roomRepository.FindByID(models.GeneralRoomID)
	if err != nil {
		generalRoom := models.Room{
			Name:     models.GeneralRoomName,
			OwnerID:  models.GeneralRoomID,
			RoomType: models.PublicRoom,
		}
		generalRoom.ID = models.GeneralRoomID

		err = s.roomRepository.Create(&generalRoom)
		if err != nil {
			log.Error().Err(err).Msg("Error creating general room")
		}
	}

	_, err = s.roomRepository.FindByID(models.AnnouncementsRoomID)
	if err != nil {
		announcementsRoom := models.Room{
			Name:     models.AnnouncementsRoomName,
			OwnerID:  models.AnnouncementsRoomID,
			RoomType: models.ReadOnlyRoom,
		}
		announcementsRoom.ID = models.AnnouncementsRoomID

		err = s.roomRepository.Create(&announcementsRoom)
		if err != nil {
			log.Error().Err(err).Msg("Error creating announcements room")
		}
	}
}
