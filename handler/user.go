package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/database/dbHelper"
	_ "github.com/sudo-abhinav/rms/middlwares"
	middlewares "github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
	"github.com/sudo-abhinav/rms/utils"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

func Createuser(w http.ResponseWriter, r *http.Request) {

	var body models.UsersRequest

	userCTX := middlewares.UserContext(r)
	createdBy := userCTX.UserID
	role := models.RoleUser

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondWithError(w, http.StatusConflict, nil, "user already exists")
		return
	}
	password, err := utils.HashPassword(body.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err, "error in hasing Password")
		return
	}

	if txErr := database.Tx(func(tx *sqlx.Tx) error {
		userId, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, password, createdBy, role)
		if saveErr != nil {
			return saveErr
		}
		//TODO :-  add logger
		log.Printf("Inserting addresses for userID %s: %+v", userId, body.Address)
		return dbHelper.CreateUserAddress(tx, userId, body.Address)
	}); txErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, "User Created...")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	userID, role, userErr := dbHelper.Login(body)
	if userErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, userErr, "failed to find user")
		return
	}

	if userID == "" || role == "" {
		utils.RespondWithError(w, http.StatusOK, nil, "user not found")
		return
	}

	sessionID, Err := dbHelper.CreateUserSession(userID)
	if Err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, Err, "failed to create user session")
		return
	}

	token, genErr := utils.GenerateJWT(userID, sessionID, role)
	if genErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, genErr, "failed to generate token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{"login successful", token})
}

func GetAllUsersByAdmin(w http.ResponseWriter, _ *http.Request) {
	users, Err := dbHelper.GetAllUser()

	if Err != nil {

		utils.RespondWithError(w, http.StatusInternalServerError, Err, "failed to get users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	sessionID := userCtx.SessionID

	if Err := dbHelper.DeleteUserSession(sessionID); Err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, Err, "failed to delete user session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"logout successful"})
}

func FetchUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCTX := middlewares.UserContext(r)
	createdBy := userCTX.UserID
	users, Err := dbHelper.FetchUserFilterBySubAdmin(createdBy)

	if Err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, Err, "failed to get users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func CalculateDistance(w http.ResponseWriter, r *http.Request) {
	var body models.DistanceRequest

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	if body.RestaurantAddressID == " " && body.UserAddressID == " " {
		utils.RespondJSON(w, http.StatusBadRequest, "Blank Data should not Accepted...")
		return
	}
	var eg errgroup.Group
	var err error
	var userCoordinates, restaurantCoordinates models.Coordinates

	eg.Go(func() error {
		userCoordinates, err = dbHelper.GetUserCoordinates(body.UserAddressID)
		return err
	})
	eg.Go(func() error {
		restaurantCoordinates, err = dbHelper.GetRestaurantCoordinates(body.RestaurantAddressID)
		return err
	})

	ergErr := eg.Wait()
	if ergErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, ergErr, "failed to get coordinates")
		return
	}

	distance, calErr := dbHelper.CalculateDistance(userCoordinates, restaurantCoordinates)
	if calErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, calErr, "failed to calculate distance")
		return
	}

	utils.RespondJSON(w, http.StatusOK, distance)
}
