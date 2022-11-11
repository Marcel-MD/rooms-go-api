package models

const (
	CreateMessage = "CreateMessage"
	UpdateMessage = "UpdateMessage"
	DeleteMessage = "DeleteMessage"
	RemoveUser    = "RemoveUser"
	AddUser       = "AddUser"
	CreateRoom    = "CreateRoom"
	UpdateRoom    = "UpdateRoom"
	DeleteRoom    = "DeleteRoom"
	Error         = "Error"
)

type Message struct {
	Base
	RoomID   string `json:"roomId"`
	Room     Room   `json:"-" gorm:"foreignKey:RoomID"`
	UserID   string `json:"userId"`
	User     User   `json:"-" gorm:"foreignKey:UserID"`
	Text     string `json:"text"`
	Command  string `json:"command"`
	TargetID string `json:"targetId"`
}
