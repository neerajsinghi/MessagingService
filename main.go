package main

import (
	"T/MessagingService/router"
	"log"
	"net/http"

	"github.com/rs/cors"
)

// setupGlobalMiddleware will setup CORS
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}

// our main function
func main() {
	router := router.NewRouter()
	log.Fatal(http.ListenAndServe(":2000", setupGlobalMiddleware(router)))

}
