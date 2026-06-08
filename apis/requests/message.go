package requests

// Message requests
type CreateMessage struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type UpdateMessage struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}
