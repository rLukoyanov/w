package requests

// Server requests
type CreateServer struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateServer struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}
