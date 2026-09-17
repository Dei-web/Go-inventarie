package types

type Users struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type ResponseData struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersCreate struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UsersUpdate struct {
	Name     string `json:"name"     validate:"omitempty,min=2,max=100"`
	Email    string `json:"email"    validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
}
