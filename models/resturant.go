package models

type CreateRestaurant struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"longitude"`
}

type Restaurant struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Address   string `json:"address" db:"address"`
	Latitude  string `json:"latitude" db:"latitude"`
	Longitude string `json:"longitude" db:"longitude"`
	CreatedBy string `json:"createdBy" db:"created_by"`
}
