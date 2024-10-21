package models

type CreateDishRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Dish struct {
	ID           string `json:"ID" db:"id"`
	Name         string `json:"name" db:"name"`
	Price        int    `json:"price" db:"price"`
	RestaurantID string `json:"resturantID" db:"restaurant_id"`
}

type RestaurantDishes struct {
	ID        string           `json:"id" db:"id"`
	Name      string           `json:"name" db:"name"`
	Address   string           `json:"address" db:"address"`
	Latitude  string           `json:"latitude" db:"latitude"`
	Longitude string           `json:"longitude" db:"longitude"`
	Dishes    []DishCollection `json:"dishes"` // Assuming you want to store multiple dishes
}

type DishCollection struct {
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}
