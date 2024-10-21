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

func FetchUserFilterBySubAdmin(created_by string) ([]models.User, error) {

	query := `select id ,name , email ,role from users where role = 'user' and	 created_by = $1 and 
                                               archived_at IS NULL`

	users := make([]models.User, 0)
	err := database.DBconn.Select(&users, query, created_by)
	if err != nil {
		return nil, err
	}

	addressQuery := `SELECT user_id, id, address, latitude, longitude FROM address
                                                WHERE user_id IN
                                                      (SELECT id FROM users WHERE archived_at IS NULL)`

	addresses := make([]models.Address, 0)
	err = database.DBconn.Select(&addresses, addressQuery)
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
