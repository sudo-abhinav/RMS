package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/database/dbHelper"
	"github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
	"github.com/sudo-abhinav/rms/utils"
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
	}

	if txErr := database.Tx(func(tx *sqlx.Tx) error {
		userId, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, password, createdBy, role)
		if saveErr != nil {
			return saveErr
		}
		return dbHelper.CreateUserAddress(tx, userId, body.Address)
	}); txErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}
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
