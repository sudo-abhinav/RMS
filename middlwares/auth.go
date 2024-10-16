package middlewares

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/sudo-abhinav/rms/database/dbHelper"
	"github.com/sudo-abhinav/rms/models"
	"github.com/sudo-abhinav/rms/utils"
	"net/http"
	"os"
	"strings"
)

type ContextKeys string

const (
	userContext ContextKeys = "userContext"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, nil, "authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.RespondWithError(w, http.StatusUnauthorized, nil, "bearer token missing")
			return
		}

		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method") // Invalid signing method error
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if parseErr != nil || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, parseErr, "invalid claims")
			return
		}

		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, nil, "invalid token ")
			return
		}

		sessionID := claimValues["sessionID"].(string)

		archivedAt, err := dbHelper.GetArchivedAt(sessionID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err, "internal server error")
			return
		}

		if archivedAt != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, nil, "invalid token")
			return
		}

		user := &models.UserCtx{
			UserID:    claimValues["userID"].(string),
			SessionID: sessionID,
			Role:      models.Role(claimValues["role"].(string)),
		}

		ctx := context.WithValue(r.Context(), userContext, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func UserContext(r *http.Request) *models.UserCtx {
	if user, ok := r.Context().Value(userContext).(*models.UserCtx); ok {
		return user
	}
	return nil
}
