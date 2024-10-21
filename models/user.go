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
	ID        string  `json:"id"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type User struct {
	ID      string    `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	Email   string    `json:"email" db:"email"`
	Address []Address `json:"address" db:"address"`
	Role    Role      `json:"role" db:"role"`
}
type Address struct {
	ID        string  `json:"id" db:"id"`
	Address   string  `json:"address" db:"address"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	UserId    string  `json:"userId" db:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginData struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password"`
	Role         Role   `db:"role"`
}
type DistanceRequest struct {
	UserAddressID       string `json:"userAddressId"`
	RestaurantAddressID string `json:"restaurantAddressId" `
}
type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Role      Role   `json:"role"`
}
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
