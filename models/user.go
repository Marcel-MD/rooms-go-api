package models

type User struct {
	Base
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"-"`

	Rooms []Room `json:"rooms" gorm:"many2many:room_users;constraint:OnDelete:CASCADE"`
}
