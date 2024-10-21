package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/rms/database/dbHelper"
	middlewares "github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
	"github.com/sudo-abhinav/rms/utils"
	"net/http"
)

func CreateDish(w http.ResponseWriter, r *http.Request) {

	restaurantID := chi.URLParam(r, "restaurantId")
	var body models.CreateDishRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, parseErr, "invalid Payload")
	}

	exist, existErr := dbHelper.IsDishExist(body.Name, restaurantID)
	if existErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, existErr, "failed to Fetch Dishes")
	}
	if exist {
		utils.RespondWithError(w, http.StatusConflict, nil, "Dish Already Exist..")
	}

	if Err := dbHelper.CreateDish(body, restaurantID); Err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, Err, "Failed to create dish")
	}
	utils.RespondJSON(w, http.StatusCreated, `Dish created...`)
}

func GetAllDish(w http.ResponseWriter, r *http.Request) {

	ListOfDishes, err := dbHelper.GetAllDish()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "failed to fetch Dishes..")
	}
	utils.RespondJSON(w, http.StatusOK, ListOfDishes)
}

func FetchAllDishesFilterBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	UserID := userCtx.UserID

	dishes, getErr := dbHelper.GetAllDishesBySubAdmin(UserID)
	if getErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, getErr, "failed to get dishes")
		return
	}

	utils.RespondJSON(w, http.StatusOK, dishes)

}

func DishesByRestaurant(w http.ResponseWriter, r *http.Request) {
	body := struct {
		RestaurantName string `json:"restaurant_name" db:"name" `
	}{}
	//w.WriteHeader(404)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	dishes, getErr := dbHelper.DishesByRestaurant(body.RestaurantName)
	if getErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, getErr, "failed to get dishes")
		return
	}

	utils.RespondJSON(w, http.StatusOK, dishes)
}
