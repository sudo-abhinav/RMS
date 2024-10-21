package dbHelper

import (
	"fmt"
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/models"
)

func CreateDish(body models.CreateDishRequest, restaurantID string) error {

	query := `INSERT INTO dishes (name, 
                    price, 
                    restaurant_id) values (trim(lower($1)) , $2 , $3 )`
	var dishID string
	Err := database.DBconn.Get(&dishID, query, body.Name, body.Price, restaurantID)
	return Err
}

func IsDishExist(name, restaurantID string) (bool, error) {
	Query :=
		`SELECT count(id) > 0 as is_exist FROM dishes
                 WHERE name = trim($1)
                   AND restaurant_id  =$2
                   AND archived_at IS NULL `

	var checkDishes bool
	err := database.DBconn.Get(&checkDishes, Query, name, restaurantID)
	return checkDishes, err
}

func GetAllDish() ([]models.Dish, error) {
	query := `SELECT id,  name ,price ,restaurant_id FROM dishes WHERE archived_at IS NULL`

	dish := make([]models.Dish, 0)
	err := database.DBconn.Select(&dish, query)
	if err != nil {
		return nil, err
	}
	return dish, nil
}

func GetAllDishesBySubAdmin(userId string) ([]models.Dish, error) {
	query := `SELECT d.id , d.name , d.price , d.restaurant_id  from dishes 
    					d INNER JOIN public.restaurants r on r.id = d.restaurant_id where created_by = $1`

	dish := make([]models.Dish, 0)
	err := database.DBconn.Select(&dish, query, userId)
	if err != nil {
		return nil, err
	}
	return dish, nil
}

//func DishesByRestaurant(name string) ([]models.RestaurantDishes, error) {
//	var restaurantDishes []models.RestaurantDishes
//
//	query := `SELECT
//					r.name ,
//					r.address,
//					r.latitude,
//					r.longitude,
//					d.id,
//					d.name,
//					d.price
//				FROM
//					public.restaurants r
//				INNER JOIN
//					dishes d ON r.id = d.restaurant_id
//				WHERE
//					r.name ILIKE $1 AND r.archived_at IS NULL AND d.archived_at IS NULL;`
//
//	err := database.DBconn.Select(&restaurantDishes, query, name)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(restaurantDishes) == 0 {
//		return nil, fmt.Errorf("restaurant not found or no dishes available")
//	}
//
//	return restaurantDishes, nil
//}

func DishesByRestaurant(name string) (models.RestaurantDishes, error) {

	// Fetch restaurant details
	restaurantQuery := `SELECT
                            r.id,
                            r.name,
                            r.address,
                            r.latitude,
                            r.longitude
                        FROM
                            public.restaurants r
                        WHERE
                            r.name ILIKE $1 AND r.archived_at IS NULL;`
	var restaurant models.RestaurantDishes
	err := database.DBconn.Get(&restaurant, restaurantQuery, name)
	if err != nil {
		return restaurant, err
	}

	// Fetch dishes for that restaurant
	dishesQuery := `SELECT name , price FROM dishes 
                    WHERE restaurant_id = $1 
                      AND archived_at IS NULL;`

	var dishes []models.DishCollection
	err = database.DBconn.Select(&dishes, dishesQuery, restaurant.ID)
	fmt.Println(dishes)
	if err != nil {
		return restaurant, err
	}

	restaurant.Dishes = dishes

	if len(restaurant.Dishes) == 0 {
		return restaurant, fmt.Errorf("restaurant found but no dishes available")
	}

	return restaurant, nil
}
