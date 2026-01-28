package server

import (
	"github.com/go-chi/chi/v5"
)

func (server *Server) userRoute(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		// All React Admin routes require authentication
		r.Use(server.authMiddleware)
		r.Delete("/", server.userController.DeleteMany)
		r.Get("/", server.userController.GetList)       // getList or getMany based on filter
		r.Post("/", server.userController.Create)       // create
		r.Get("/{id}", server.userController.GetOne)    // getOne
		r.Put("/{id}", server.userController.Update)    // update
		r.Delete("/{id}", server.userController.Delete) // delete
	})
}
