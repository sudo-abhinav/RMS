package routes

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sudo-abhinav/rms/handler"
	"github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/routes/admin"
	sub_admin "github.com/sudo-abhinav/rms/routes/sub-admin"
	"github.com/sudo-abhinav/rms/routes/users"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	commonMiddlewares := middlewares.CommonMiddlewares()
	router.Use(commonMiddlewares...)
	router.Route("/test", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("this all about RMS")
		})
	})
	router.Route("/v1", func(r chi.Router) {

		r.Post("/login", handler.LoginUser)

		r.Group(func(r chi.Router) {

			r.Use(middlewares.Authenticate)
			r.Use(middleware.Logger)

			/* TODO error in this route
			check how to handle data
			*/
			r.Get("/restaurantDishes", handler.DishesByRestaurant)

			r.Post("/logout", handler.LogoutUser)

			//Admin route
			admin.AdminRoutes(r)
			// Sub-admin routes
			sub_admin.SubAdminRoutes(r)
			//	user route
			/*TODO = not completed yet */
			users.UsersRoutes(r)

		})
	})

	return &Server{
		Router: router,
	}
}

func (server *Server) RUN(PORT string) error {
	server.server = &http.Server{
		Addr:              PORT,
		Handler:           server.Router,
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      5 * time.Minute,
	}
	return server.server.ListenAndServe()

}

func (server *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return server.server.Shutdown(ctx)
}
