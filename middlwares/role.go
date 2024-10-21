package middlewares

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/rms/database/dbHelper"
	"github.com/sudo-abhinav/rms/models"
	"net/http"
)

func ShouldHaveRole(role models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := UserContext(r).Role
			if userRole != role {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CheckSubAdminRestaurant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		restaurantID := chi.URLParam(r, "restaurantId")
		userCTX := UserContext(r).UserID
		createdBy := userCTX

		// Logic to check if the restaurant was created by the sub-admin
		exist, err := dbHelper.IsRestaurantCreatedBySubAdmin(restaurantID, createdBy)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !exist {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}
