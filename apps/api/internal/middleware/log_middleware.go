package middleware

import (
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		// You can customize the logging format as needed
		log.Printf("Received %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		log.Printf("Query Parameters : %s", r.URL.RawQuery)
		log.Printf("Headers: %v", r.Header)
		log.Printf("Request Body : %s", r.Body)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log after the response is sent
		log.Printf("Completed %s request for %s", r.Method, r.URL.Path)
	})

}
