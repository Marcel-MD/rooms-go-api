package models

type AllowedUser struct {
	Email       string `json:"email" gorm:"primaryKey"`
	DefaultRole string `json:"defaultRole"`
}
