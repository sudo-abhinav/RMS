package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/rms/handler"
	middlewares "github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
)

func UsersRoutes(r chi.Router) {
	r.Route("/user", func(user chi.Router) {
		user.Use(middlewares.ShouldHaveRole(models.RoleSubAdmin))

		user.Get("/allRestaurant", handler.GetAllRestaurant)
		user.Get("/allDishes", handler.GetAllDish)

		user.Get("/calculate-Distance", handler.CalculateDistance)

	})
}
