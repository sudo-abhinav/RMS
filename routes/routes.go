package routes

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sudo-abhinav/rms/handler"
	"github.com/sudo-abhinav/rms/middlwares"
	"github.com/sudo-abhinav/rms/models"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Route("/test", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("this all about RMS")
		})
	})
	router.Route("/v1", func(r chi.Router) {

		r.Post("/login", handler.LoginUser)

		router.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate)

			//r.Get("/dishes-by-restaurant", handler.DishesByRestaurant)
			//r.Post("/logout", handler.LogoutUser)

			r.Route("/admin", func(admin chi.Router) {
				r.Use(middlewares.ShouldHaveRole(models.RoleAdmin))

				admin.Post("/create-user", handler.Createuser)
				admin.Get("/all-users", handler.GetAllUsersByAdmin)

				admin.Post("/create-SubAdmin", handler.CreateSubAdmin)
				admin.Get("/all-subAdmin", handler.SeeAllSUbAdmin)
			})
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
