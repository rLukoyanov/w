package requests

// Channel requests
type CreateChannel struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Type string `json:"type" validate:"required,oneof=text voice"`
}
