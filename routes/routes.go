package routes

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("this all about RMS")
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
