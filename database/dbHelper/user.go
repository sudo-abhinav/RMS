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
	Err := tx.Get(&userID, SQL, name, email, password, createdBy, role)
	return userID, Err
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
	query = utils.SetupBindVars(query, "(? , ? , ? , ? )", len(addresses))
	_, err := tx.Exec(query, data...)
	return err
}

func GetAllUser() ([]models.User, error) {
	usersQuery := `SELECT id, name, email, role FROM users WHERE archived_at IS NULL ORDER BY id`
	// it creates a slice of models.User with an initial length of 0.
	users := make([]models.User, 0)
	err := database.DBconn.Select(&users, usersQuery)
	if err != nil {
		return nil, err
	}

	addressesQuery := `SELECT user_id, id, address, latitude, longitude FROM address 
                                                 WHERE user_id IN 
                                                       (SELECT id FROM users
                                                                  WHERE archived_at IS NULL)`
	addresses := make([]models.Address, 0)
	err = database.DBconn.Select(&addresses, addressesQuery)
	if err != nil {
		return nil, err
	}
	addressMap := make(map[string][]models.Address)
	for _, address := range addresses {
		addressMap[address.UserId] = append(addressMap[address.UserId], address)
	}
	for i := range users {
		users[i].Address = addressMap[users[i].ID]
	}

	return users, nil
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

func Login(body models.LoginRequest) (string, models.Role, error) {
	SQL := `SELECT u.id,
       			   u.role,
       			   u.password
			  FROM users u
			  WHERE u.email = TRIM($1)
			    AND u.archived_at IS NULL`

	var user models.LoginData
	if getErr := database.DBconn.Get(&user, SQL, body.Email); getErr != nil {
		return "", "", getErr
	}
	if passwordErr := utils.VerifyPassword(body.Password, user.PasswordHash); passwordErr != nil {
		return "", "", passwordErr
	}
	return user.ID, user.Role, nil
}

func CreateUserSession(userID string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session(user_id)
    			VALUES ($1) RETURNING id `

	Err := database.DBconn.Get(&sessionID, query, userID)
	if Err != nil {
		return "", Err
	}
	return sessionID, nil
}

func DeleteUserSession(sessionID string) error {

	query := `UPDATE user_session SET 
                        archived_at = now() 
                    where user_id = $1 AND 
                          archived_at is NULL`

	_, err := database.DBconn.Exec(query, sessionID)
	return err
}
