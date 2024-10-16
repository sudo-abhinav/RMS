package models

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleSubAdmin Role = "sub-admin"
	RoleUser     Role = "user"
)

type UsersRequest struct {
	ID       string           `json:"id"`
	Name     string           `json:"name" `
	Email    string           `json:"email"`
	Password string           `json:"password"`
	Address  []AddressRequest `json:"address"`
}
type AddressRequest struct {
	ID        string `json:"id"`
	Address   string
	Latitude  string
	Longitude string
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Role      Role   `json:"role"`
}
