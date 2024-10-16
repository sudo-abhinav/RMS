package utils

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// w http.ResponseWriter, statusCode int, body interface{}
func RespondWithError(w http.ResponseWriter, code int, err error, body interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		logrus.Errorf("Failed to send error to caller with error: %+v", err)
	}
}

// ParseBody parses the values from io reader to a given interface
func ParseBody(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}

	return nil
}

// EncodeJSONBody writes the JSON body to response writer
func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

// RespondJSON sends the interface as a JSON
func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("Failed to respond JSON : %+v", err)
		}
	}
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Verify Password verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}

// GenerateJWT creates a new JWT token with claims
func GenerateJWT(userID, email, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"userID":    userID,
		"email":     email,
		"sessionID": sessionID,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 3).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func SetupBindVars(stmt, bindVars string, length int) string {
	bindVars += ","
	stmt = fmt.Sprintf("%s %s", stmt, strings.Repeat(bindVars, length))
	return replaceSQL(strings.TrimSuffix(stmt, ","), "?")
}

func replaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
