package models

type Room struct {
	Base
	Name     string    `json:"name"`
	OwnerID  string    `json:"ownerId"`
	Users    []User    `json:"users" gorm:"many2many:room_users;constraint:OnDelete:CASCADE"`
	Messages []Message `json:"-" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
	RoomType string    `json:"roomType"`
}

const (
	PrivateRoom  = "Private"
	PublicRoom   = "Public"
	ReadOnlyRoom = "ReadOnly"
)

const (
	AnnouncementsRoomID   = "announcements"
	AnnouncementsRoomName = "Announcements"

	GeneralRoomID   = "general"
	GeneralRoomName = "General"
)
