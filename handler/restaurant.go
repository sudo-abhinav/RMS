package handler

import (
	"github.com/sudo-abhinav/rms/database/dbHelper"
	middlewares "github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
	"github.com/sudo-abhinav/rms/utils"
	"net/http"
)

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var RestaurantData models.CreateRestaurant
	userCtx := middlewares.UserContext(r)
	createdBy := userCtx.UserID

	if parseErr := utils.ParseBody(r.Body, &RestaurantData); parseErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	exists, existsErr := dbHelper.IsRestatuarntExist(RestaurantData.Name, RestaurantData.Address)
	if existsErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, existsErr, "Failed to check Restaurant Existence")
		return
	}
	if exists {
		utils.RespondWithError(w, http.StatusConflict, nil, "Restaurant Already Exist..")
	}
	if saveErr := dbHelper.CreateRestaurant(RestaurantData, createdBy); saveErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, saveErr, "Error In Register Restaurant..")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, "Restaurant Registered")
}
func GetAllRestaurant(w http.ResponseWriter, r *http.Request) {

	ListofRestaurant, Err := dbHelper.GetAllRestaurant()

	if Err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, Err, "failed to Fetch Restaurant")
		return
	}

	utils.RespondJSON(w, http.StatusOK, ListofRestaurant)
}

func GetRestaurantCreatedBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	RestaurantData, err := dbHelper.RestaurantCreatedBySubAdmin(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err, "error in Fetching Restaurants")
		return
	}

	utils.RespondJSON(w, http.StatusOK, RestaurantData)
}
