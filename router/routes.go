package router

import (
	"net/http"

	handler "T/MessagingService/clients"
)

// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	Route{
		"ws",
		"GET",
		"/ws",
		handler.ServeWebSocket,
	},
}
