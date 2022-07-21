package repositories

import (
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
	log.Debug().Msg("Finding all rooms")

	var rooms []models.Room
	r.DB.Find(&rooms)
	return rooms
}

func (r *RoomRepository) FindByID(id string) (models.Room, error) {
	log.Debug().Str("id", id).Msg("Finding room")

	var room models.Room
	err := r.DB.First(&room, "id = ?", id).Error

	return room, err
}

func (r *RoomRepository) FindByIdWithUsers(id string) (models.Room, error) {
	log.Debug().Str("id", id).Msg("Finding room")

	var room models.Room
	err := r.DB.Model(&models.Room{}).Preload("Users").First(&room, "id = ?", id).Error

	return room, err
}

func (r *RoomRepository) Create(room *models.Room) error {
	log.Debug().Msg("Creating room")

	return r.DB.Create(room).Error
}

func (r *RoomRepository) Update(room *models.Room) error {
	log.Debug().Msg("Updating room")

	return r.DB.Save(room).Error
}

func (r *RoomRepository) Delete(room *models.Room) error {
	log.Debug().Msg("Deleting room")

	return r.DB.Delete(room).Error
}

func (r *RoomRepository) AddUser(room *models.Room, user *models.User) error {
	log.Debug().Msg("Adding user to room")

	return r.DB.Model(room).Omit("Users.*").Association("Users").Append(user)
}

func (r *RoomRepository) RemoveUser(room *models.Room, user *models.User) error {
	log.Debug().Msg("Removing user from room")

	return r.DB.Model(room).Association("Users").Delete(user)
}
