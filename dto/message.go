package dto

type CreateMessage struct {
	Text string `json:"text" binding:"required,min=1,max=500"`
}

type UpdateMessage struct {
	Text string `json:"text" binding:"required,min=1,max=500"`
}

type MessageQueryParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
