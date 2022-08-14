package dto

type RegisterUser struct {
	FirstName string `json:"firstName" binding:"required,min=3,max=50"`
	LastName  string `json:"lastName" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=50"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type UpdateUser struct {
	FirstName string `json:"firstName" binding:"required,min=3,max=50"`
	LastName  string `json:"lastName" binding:"required,min=3,max=50"`
}

type SearchByEmail struct {
	Email string `json:"email" binding:"required"`
}
