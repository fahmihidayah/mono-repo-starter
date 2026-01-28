package server

import (
	"github.com/go-chi/chi/v5"
)

func (server *Server) postRoute(r chi.Router) {
	r.Route("/posts", func(r chi.Router) {
		// All React Admin routes require authentication
		r.Use(server.authMiddleware)

		r.Delete("/", server.postController.DeleteMany)
		r.Put("/", server.postController.UpdateMany)

		// Handle single resource operations
		r.Get("/", server.postController.GetList)       // getList or getMany based on filter
		r.Post("/", server.postController.Create)       // create
		r.Get("/{id}", server.postController.GetOne)    // getOne
		r.Put("/{id}", server.postController.Update)    // update
		r.Delete("/{id}", server.postController.Delete) // delete
	})
}
