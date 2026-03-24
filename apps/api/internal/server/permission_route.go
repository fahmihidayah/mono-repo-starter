package server

import (
	"github.com/go-chi/chi/v5"
)

func (server *Server) permissionRoute(r chi.Router) {
	r.Route("/permissions", func(r chi.Router) {
		// All React Admin routes require authentication
		r.Use(server.authMiddleware)

		// Handle bulk operations first (they have query parameters)
		r.Delete("/", server.permissionController.DeleteMany)
		r.Put("/", server.permissionController.UpdateMany)

		// Handle single resource operations
		r.Get("/", server.permissionController.GetList)       // getList or getMany based on filter
		r.Post("/", server.permissionController.Create)       // create
		r.Get("/{id}", server.permissionController.GetOne)    // getOne
		r.Put("/{id}", server.permissionController.Update)    // update
		r.Delete("/{id}", server.permissionController.Delete) // delete
	})
}
