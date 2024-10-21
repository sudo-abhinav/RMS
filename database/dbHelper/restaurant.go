package dbHelper

import (
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/models"
)

func IsRestatuarntExist(name, address string) (bool, error) {

	Query :=
		`SELECT count(id) > 0 as is_exist FROM restaurants
                 WHERE name = trim($1)
                   AND address  =trim($2)
                   AND archived_at IS NULL `

	var checkRestaurant bool
	err := database.DBconn.Get(&checkRestaurant, Query, name, address)
	return checkRestaurant, err
}

func CreateRestaurant(body models.CreateRestaurant, userId string) error {

	values := []interface{}{body.Name, body.Address, body.Latitude, body.Longitude, userId}
	query := `INSERT INTO restaurants (name, address, latitude, longitude, created_by)
							values (lower($1) , lower($2), $3 ,$4 , $5)`
	_, Err := database.DBconn.Exec(query, values...)
	return Err

}

func GetAllRestaurant() ([]models.Restaurant, error) {
	//language = sql
	sql := `SELECT id ,
				   name ,
				   address ,
				   latitude ,
				   longitude , 
				   created_by FROM restaurants 
				              where archived_at 
				                  IS NULL order by created_by`

	// it creates a slice of models.User with an initial length of 0.
	restaurants := make([]models.Restaurant, 0)
	err := database.DBconn.Select(&restaurants, sql)
	return restaurants, err
}

func RestaurantCreatedBySubAdmin(userID string) ([]models.Restaurant, error) {
	query := `SELECT id,
       				name , 
       				address ,
       				longitude ,
       				latitude ,
       				created_by from restaurants where created_by = '$1'`
	data := make([]models.Restaurant, 0)
	err := database.DBconn.Get(&data, query, userID)
	return data, err
}

// func IsRestaurantCreatedBySubAdmin(restaurantID string, userID string) (bool, error) {
//
//		query := `select restaurants.id , created_by from restaurants where created_by=$1 and restaurants.id = $2`
//		err := database.DBconn.Select(query, restaurantID, userID)
//		if err != nil {
//			return false, err
//		}
//		return true, err
//	}
func IsRestaurantCreatedBySubAdmin(restaurantID string, userID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM restaurants WHERE created_by = $1 AND id = $2`
	err := database.DBconn.Get(&count, query, userID, restaurantID)
	//.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
