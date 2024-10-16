package dbHelper

import (
	"github.com/jmoiron/sqlx"
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/models"
	"github.com/sudo-abhinav/rms/utils"
	"time"
)

func CreateUser(tx *sqlx.Tx, name, email, password, createdBy string, role models.Role) (string, error) {
	SQL := `INSERT INTO users (name, email, password, created_by, role)
			  VALUES (TRIM($1), TRIM($2), $3, $4, $5) RETURNING id`

	var userID string
	crtErr := tx.Get(&userID, SQL, name, email, password, createdBy, role)
	return userID, crtErr
}
func CreateUserAddress(tx *sqlx.Tx, userID string, addresses []models.AddressRequest) error {
	query := `INSERT INTO address (user_id, address, latitude, longitude) VALUES`
	
	data := make([]interface{}, 0)
	for i := range addresses {
		data = append(data,
			userID,
			addresses[i].Address,
			addresses[i].Latitude,
			addresses[i].Longitude,
		)
	}
	utils.SetupBindVars(query, " (? , ? , ? ,?)", len(addresses))
	_, err := tx.Exec(query, data...)
	return err
}

func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	SQL := `SELECT archived_at 
              FROM user_session 
              WHERE id = $1
              	AND archived_at IS NULL`

	getErr := database.DBconn.Get(&archivedAt, SQL, sessionID)
	return archivedAt, getErr
}
func IsUserExists(email string) (bool, error) {
	Query :=
		`SELECT count(id) > 0 as is_exist FROM users
                 WHERE email = trim(lower($1)) 
                   AND archived_at IS NULL `

	var checkUser bool
	err := database.DBconn.Get(&checkUser, Query, email)
	return checkUser, err
}
