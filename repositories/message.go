package repositories

import (
	"sync"

	"github.com/Marcel-MD/rooms-go-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type IMessageRepository interface {
	FindByRoomID(roomID string, page, size int) []models.Message
	FindByID(id string) (models.Message, error)
	Create(message *models.Message) error
	Update(message *models.Message) error
	Delete(message *models.Message) error
}

type MessageRepository struct {
	DB *gorm.DB
}

var (
	messageOnce       sync.Once
	messageRepository IMessageRepository
)

func GetMessageRepository() IMessageRepository {
	messageOnce.Do(func() {
		log.Info().Msg("Initializing message repository")
		messageRepository = &MessageRepository{
			DB: models.GetDB(),
		}
	})
	return messageRepository
}

func (r *MessageRepository) FindByRoomID(roomID string, page, size int) []models.Message {
	var messages []models.Message

	r.DB.Scopes(models.Paginate(page, size)).Model(&models.Message{}).
		Order("created_at desc").Preload("User").Find(&messages, "room_id = ?", roomID)

	return messages
}

func (r *MessageRepository) FindByID(id string) (models.Message, error) {
	var message models.Message
	err := r.DB.First(&message, "id = ?", id).Error

	return message, err
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.DB.Create(message).Error
}

func (r *MessageRepository) Update(message *models.Message) error {
	return r.DB.Save(message).Error
}

func (r *MessageRepository) Delete(message *models.Message) error {
	return r.DB.Delete(message).Error
}
