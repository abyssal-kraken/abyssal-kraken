package http

type CreateAccountRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
