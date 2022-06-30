package models

type Message struct {
	Base
	RoomID string `json:"room_id"`
	Room   Room   `json:"room" gorm:"foreignKey:RoomID"`
	UserID string `json:"user_id"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Text   string `json:"text"`
}
