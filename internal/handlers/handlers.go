package handlers

import (
	"github.com/dscott/invoicey/internal/services"
)

// Handlers contains all HTTP handlers with their dependencies
type Handlers struct {
	Services *services.Services
}

// New creates a new Handlers instance with all dependencies
func New(svc *services.Services) *Handlers {
	return &Handlers{
		Services: svc,
	}
}
