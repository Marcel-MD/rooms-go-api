package dto

type CreateMessage struct {
	Text string `json:"text" binding:"required,min=3,max=500"`
}

type UpdateMessage struct {
	Text string `json:"text" binding:"required,min=3,max=500"`
}

type MessageQueryParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
