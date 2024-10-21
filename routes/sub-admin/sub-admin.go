package sub_admin

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/rms/handler"
	middlewares "github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
)

func SubAdminRoutes(r chi.Router) {
	r.Route("/sub-admin", func(subAdmin chi.Router) {
		subAdmin.Use(middlewares.ShouldHaveRole(models.RoleSubAdmin))

		subAdmin.Post("/create-user", handler.Createuser)
		subAdmin.Get("/all-users", handler.FetchUsersBySubAdmin)

		subAdmin.Post("/create-restaurant", handler.CreateRestaurant)
		subAdmin.Get("/get-restaurant", handler.GetRestaurantCreatedBySubAdmin)

		subAdmin.Route("/{restaurantId}", func(r chi.Router) {
			/* TODO
			=> i have add a middleware here to check the restaurant is created by sub-admin or not
			=> if sub-admin access the restaurant created by admin this thing is bug ..
			=> extract current logged user id form context and check
			*/
			r.Use(middlewares.CheckSubAdminRestaurant)
			r.Post("/create-dish", handler.CreateDish)
		})
		subAdmin.Get("/", handler.FetchAllDishesFilterBySubAdmin)
	})
}
