package server

import (
	"github.com/go-chi/chi/v5"
)

func (server *Server) categoryRoute(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		// All React Admin routes require authentication
		r.Use(server.authMiddleware)

		// Handle bulk operations first (they have query parameters)
		r.Delete("/", server.categoryController.DeleteMany)
		r.Put("/", server.categoryController.UpdateMany)

		// Handle single resource operations
		r.Get("/", server.categoryController.GetList)       // getList or getMany based on filter
		r.Post("/", server.categoryController.Create)       // create
		r.Get("/{id}", server.categoryController.GetOne)    // getOne
		r.Put("/{id}", server.categoryController.Update)    // update
		r.Delete("/{id}", server.categoryController.Delete) // delete
	})
}
