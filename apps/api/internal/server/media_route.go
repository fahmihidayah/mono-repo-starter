package server

import (
	"github.com/go-chi/chi/v5"
)

func (server *Server) mediaRoute(r chi.Router) {
	r.Route("/media", func(r chi.Router) {
		// All React Admin routes require authentication
		r.Use(server.authMiddleware)

		// Handle bulk operations first (they have query parameters)
		// r.Method("PUT", "/", server.routeToUpdateManyMedia())
		// r.Method("DELETE", "/", server.routeToDeleteMany("media"))

		r.Delete("/", server.mediaController.DeleteMany)
		r.Put("/", server.mediaController.UpdateMany)
		// Handle single resource operations
		r.Get("/", server.mediaController.GetList)       // getList or getMany based on filter
		r.Post("/", server.mediaController.Create)       // create (file upload)
		r.Get("/{id}", server.mediaController.GetOne)    // getOne
		r.Put("/{id}", server.mediaController.Update)    // update
		r.Delete("/{id}", server.mediaController.Delete) // delete
	})
}

// func (server *Server) routeToUpdateManyMedia() chi.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if server.mediaController.HasFilterParam(r) {
// 			server.mediaController.UpdateMany(w, r)
// 		} else {
// 			server.mediaController.SendBadRequest(w, "Filter parameter required for bulk update")
// 		}
// 	}
// }
