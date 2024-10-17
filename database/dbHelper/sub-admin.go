package dbHelper

import (
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/models"
)

func CreateSubAdmin(name, email, password, createdBy string, role models.Role) error {
	SQL := `INSERT INTO users (name, email, password, created_by, role)
			  VALUES (TRIM($1), TRIM($2), $3, $4, $5) RETURNING id`

	var userID string
	crtErr := database.DBconn.Get(&userID, SQL, name, email, password, createdBy, role)
	return crtErr
}

func GetAllSubAdmins() ([]models.SubAdmin, error) {
	query := `SELECT id ,
       				name ,
       				email,
       				role ,
       				created_by 
							FROM users where role='sub-admin' AND archived_at IS NULL `

	subAdmins := make([]models.SubAdmin, 0)
	Err := database.DBconn.Select(&subAdmins, query)
	return subAdmins, Err
}
