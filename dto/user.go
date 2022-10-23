package dto

type RegisterUser struct {
	FirstName string `json:"firstName" binding:"required,min=3,max=50"`
	LastName  string `json:"lastName" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,min=8,max=10"`
	Password  string `json:"password" binding:"required,min=8,max=50"`
}

type RegisterOtpUser struct {
	RegisterUser
	Otp string `json:"otp" binding:"required,len=6"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type LoginOtpUser struct {
	LoginUser
	Otp string `json:"otp" binding:"required,len=6"`
}

type UpdateUser struct {
	FirstName string `json:"firstName" binding:"required,min=3,max=50"`
	LastName  string `json:"lastName" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,min=8,max=10"`
}

type UpdateOtpUser struct {
	UpdateUser
	Otp string `json:"otp" binding:"required,len=6"`
}

type SearchByEmail struct {
	Email string `json:"email" binding:"required"`
}

type Email struct {
	Email string `json:"email" binding:"required,email"`
}
