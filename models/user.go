package models

type User struct {
	Base
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"-"`

	Rooms []Room `json:"rooms" gorm:"many2many:room_users;constraint:OnDelete:CASCADE"`
}
