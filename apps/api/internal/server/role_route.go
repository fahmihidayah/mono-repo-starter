package server

import (
	"github.com/go-chi/chi/v5"
)

func (server *Server) roleRoute(r chi.Router) {
	r.Route("/roles", func(r chi.Router) {
		// All React Admin routes require authentication
		r.Use(server.authMiddleware)

		// Handle bulk operations first (they have query parameters)
		r.Delete("/", server.roleController.DeleteMany)
		r.Put("/", server.roleController.UpdateMany)

		// Handle single resource operations
		r.Get("/", server.roleController.GetList)       // getList or getMany based on filter
		r.Post("/", server.roleController.Create)       // create
		r.Get("/{id}", server.roleController.GetOne)    // getOne
		r.Put("/{id}", server.roleController.Update)    // update
		r.Delete("/{id}", server.roleController.Delete) // delete
	})
}
