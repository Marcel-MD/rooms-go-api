package models

type Room struct {
	Base
	Name     string    `json:"name"`
	OwnerID  string    `json:"owner_id"`
	Users    []User    `json:"users" gorm:"many2many:room_users;constraint:OnDelete:CASCADE"`
	Messages []Message `json:"-" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
}
