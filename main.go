package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"T/MessagingService/router"

)

// setupGlobalMiddleware will setup CORS
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}

// our main function
func main() {

	router := router.NewRouter()

	log.Fatal(http.ListenAndServe(":6008", setupGlobalMiddleware(router)))

}
