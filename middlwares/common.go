package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
)

// corsOptions setting up routes for cors
func corsOptions() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Token", "importDate", "X-Client-Version", "Cache-Control", "Pragma", "x-started-at", "x-api-key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})
}

// CommonMiddlewares middleware common for all routes
//func CommonMiddlewares() chi.Middlewares {
//	return chi.Chain(
//		func(next http.Handler) http.Handler {
//			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				w.Header().Add("Content-Type", "application/json")
//
//				// request
//				next.ServeHTTP(w, r)
//				// response
//			})
//		},
//		corsOptions().Handler,
//		func(next http.Handler) http.Handler {
//			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				defer func() {
//					err := recover()
//					if err != nil {
//						log.Println("common error")
//						//log.Logger.Errorf("Request Panic err: %v", err)
//						jsonBody, _ := json.Marshal(map[string]string{
//							"error": "There was an internal server error",
//						})
//						w.Header().Set("Content-Type", "application/json")
//						w.WriteHeader(http.StatusInternalServerError)
//						_, err := w.Write(jsonBody)
//						if err != nil {
//							log.Println("common middleware error 2")
//							//log.Logger.Errorf("Failed to send response from middleware with error: %+v", err)
//						}
//					}
//				}()
//				next.ServeHTTP(w, r)
//			})
//		},
//	)
//}

func CommonMiddlewares() chi.Middlewares {
	return chi.Middlewares{
		// Default headers
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			})
		},
		// CORS handling
		corsOptions().Handler,
		// request and response logging middle-ware for securiyy purpose
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				req := map[string]interface{}{
					"header":      r.Header.Clone(),
					"method":      r.Method,
					"url":         r.URL,
					"request-uri": r.RequestURI,
					"params":      r.URL.Query(),
					"client-ip":   r.RemoteAddr,
				}

				logrus.WithContext(r.Context()).WithField("request", req).Info("Request received")

				rec := httptest.NewRecorder()
				next.ServeHTTP(rec, r)
				resp := map[string]interface{}{
					"status-code": rec.Code,
					"response":    rec.Body.String(),
					"header":      rec.Header(),
				}
				logrus.WithContext(r.Context()).WithField("response", resp).Info("Request completed")

				for k, v := range rec.Header() {
					w.Header()[k] = v
				}
				w.WriteHeader(rec.Code)
				_, err := rec.Body.WriteTo(w)
				if err != nil {
					logrus.WithContext(r.Context()).WithError(err).Error("Failed to write response")
				}
			})
		},
		// Panic recovery middleware
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						logrus.Errorf("Request Panic: %v", err)
						jsonBody, _ := json.Marshal(map[string]interface{}{
							"error": "Internal server error",
							"trace": fmt.Sprintf("%+v", err),
							"stack": string(debug.Stack()),
						})
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_, writeErr := w.Write(jsonBody)
						if writeErr != nil {
							logrus.Errorf("Failed to send response: %v", writeErr)
						}
					}
				}()
				next.ServeHTTP(w, r)
			})
		},
	}
}
