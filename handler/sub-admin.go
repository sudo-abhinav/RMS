package handler

import (
	"github.com/sudo-abhinav/rms/database/dbHelper"
	middlewares "github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
	"github.com/sudo-abhinav/rms/utils"
	"net/http"
)

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	var body models.SubAdminRequest

	userCTX := middlewares.UserContext(r)
	createdBy := userCTX.UserID
	role := models.RoleSubAdmin

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
	hashPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err, "error in hasing Password")
	}

	if saveErr := dbHelper.CreateSubAdmin(body.Name, body.Email, hashPassword, createdBy, role); saveErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, saveErr, "failed to create sub-admin")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, "sub_admin created successfully")
}

// TODO :- `_` is used here to ignore the second return value from the http ListenAndServe function
func SeeAllSUbAdmin(w http.ResponseWriter, _ *http.Request) {
	allSubsdmin, Err := dbHelper.GetAllSubAdmins()
	if Err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, Err, "error in Fetching Sub-admin list")
	}

	utils.RespondJSON(w, http.StatusOK, allSubsdmin)
}
