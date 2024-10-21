package admin

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/rms/handler"
	_ "github.com/sudo-abhinav/rms/handler"
	middlewares "github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
)

func AdminRoutes(r chi.Router) {
	r.Route("/admin", func(admin chi.Router) {
		admin.Use(middlewares.ShouldHaveRole(models.RoleAdmin))

		admin.Post("/create-user", handler.Createuser)
		admin.Get("/all-users", handler.GetAllUsersByAdmin)

		admin.Post("/create-SubAdmin", handler.CreateSubAdmin)
		admin.Get("/all-subAdmin", handler.SeeAllSUbAdmin)

		admin.Post("/create-restaurant", handler.CreateRestaurant)
		admin.Get("/all-restaurant", handler.GetAllRestaurant)

		// Restaurant-specific routes
		admin.Route("/{restaurantId}", func(r chi.Router) {
			r.Post("/create-dish", handler.CreateDish)
		})
		admin.Get("/all-dish", handler.GetAllDish)
	})
}
