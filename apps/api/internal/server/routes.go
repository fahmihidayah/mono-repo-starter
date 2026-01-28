package server

import (
	"net/http"

	"github.com/fahmihidayah/go-api-orchestrator/internal/handler"
	appMiddleware "github.com/fahmihidayah/go-api-orchestrator/internal/middleware"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (server *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(appMiddleware.CorsMiddleware(server.config))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.SendSuccess(w, "Ok", nil)
	})

	r.NotFound(handler.NotFoundHandler)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	// Serve static files from uploads directory (for local storage)
	// This will serve files at http://localhost:8080/uploads/filename.jpg
	if server.config.Storage.Type == "local" {
		uploadDir := server.config.Storage.LocalUploadDir
		fileServer := http.FileServer(http.Dir(uploadDir))
		r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))
	}

	r.Route("/api", func(r chi.Router) {
		server.userRoute(r)
		server.authRoute(r)
		server.postRoute(r)
		server.mediaRoute(r)
		server.categoryRoute(r)
	})

	return r
}
