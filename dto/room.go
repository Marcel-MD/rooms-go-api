package dto

type CreateRoom struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

type UpdateRoom struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}
