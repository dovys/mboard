package handlers

import (
	"github.com/gorilla/mux"
)

// Handler provides an interface to register routes
type Handler interface {
	Register(*mux.Router)
}
