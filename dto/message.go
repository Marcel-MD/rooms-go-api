package dto

type WebSocketMessage struct {
	Text     string `json:"text" binding:"required,max=500"`
	Command  string `json:"command" binding:"required,min=1,max=50"`
	TargetID string `json:"targetId" binding:"required,min=1,max=50"`
	RoomID   string `json:"roomId" binding:"required,min=1,max=50"`
}

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
