package repositories

import (
	"errors"
	"sync"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type IRoomRepository interface {
	FindAll() []models.Room
	FindByID(id string) (models.Room, error)
	FindByIdWithUsers(id string) (models.Room, error)
	Create(room *models.Room) error
	Update(room *models.Room) error
	Delete(room *models.Room) error
	AddUser(room *models.Room, user *models.User) error
	RemoveUser(room *models.Room, user *models.User) error
	VerifyUserInRoom(roomID, userID string) error
}

type RoomRepository struct {
	DB *gorm.DB
}

var (
	roomOnce       sync.Once
	roomRepository IRoomRepository
)

func GetRoomRepository() IRoomRepository {
	roomOnce.Do(func() {
		log.Info().Msg("Initializing room repository")
		roomRepository = &RoomRepository{
			DB: models.GetDB(),
		}
	})
	return roomRepository
}

func (r *RoomRepository) FindAll() []models.Room {
	var rooms []models.Room
	r.DB.Find(&rooms)
	return rooms
}

func (r *RoomRepository) FindByID(id string) (models.Room, error) {
	var room models.Room
	err := r.DB.First(&room, "id = ?", id).Error

	return room, err
}

func (r *RoomRepository) FindByIdWithUsers(id string) (models.Room, error) {
	var room models.Room
	err := r.DB.Model(&models.Room{}).Preload("Users").First(&room, "id = ?", id).Error

	return room, err
}

func (r *RoomRepository) Create(room *models.Room) error {
	return r.DB.Create(room).Error
}

func (r *RoomRepository) Update(room *models.Room) error {
	return r.DB.Save(room).Error
}

func (r *RoomRepository) Delete(room *models.Room) error {
	return r.DB.Delete(room).Error
}

func (r *RoomRepository) AddUser(room *models.Room, user *models.User) error {
	return r.DB.Model(room).Omit("Users.*").Association("Users").Append(user)
}

func (r *RoomRepository) RemoveUser(room *models.Room, user *models.User) error {
	return r.DB.Model(room).Association("Users").Delete(user)
}

func (r *RoomRepository) VerifyUserInRoom(roomID, userID string) error {
	var room models.Room
	err := r.DB.Model(&models.Room{}).Preload("Users").First(&room, "id = ?", roomID).Error
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
