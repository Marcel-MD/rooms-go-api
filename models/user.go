package models

import "github.com/lib/pq"

type User struct {
	Base
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Phone     string `json:"-"`
	Password  string `json:"-"`

	Roles    pq.StringArray `json:"roles" gorm:"type:text[]"`
	IsOnline bool           `json:"isOnline"`

	Rooms []Room `json:"rooms" gorm:"many2many:room_users;constraint:OnDelete:CASCADE"`
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}

const (
	UserRole  = "user"
	AdminRole = "admin"
)
