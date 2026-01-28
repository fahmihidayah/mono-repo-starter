package server

import "github.com/go-chi/chi/v5"

func (server *Server) authRoute(r chi.Router) {
	r.Route("/users/auth", func(r chi.Router) {
		// Public routes
		r.Post("/register", server.authController.RegisterUser)
		r.Post("/verify", server.authController.VerifyUser)
		r.Post("/login", server.authController.Login)
		r.Post("/initial-reset-password", server.authController.InitialResetPassword)
		r.Post("/complete-reset-password", server.authController.CompleteResetPassword)

		// Protected routes - require authentication
		r.Group(func(r chi.Router) {
			r.Use(server.authMiddleware)
			r.Post("/logout", server.authController.Logout)
		})
	})
}
