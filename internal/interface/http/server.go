package http

import (
	"github.com/aidapedia/gdk/http/server"
)

// HTTPServiceInterface is an interface to handle http service
type HTTPServiceInterface interface {
	ListenAndServe() error
	// Stop()
}

// HTTPService is a struct to handle http service
type HTTPService struct {
	svr *server.Server
}

// NewHTTPService is a function to create a new http service
func NewHTTPService() HTTPServiceInterface {
	svr, _ := server.NewWithDefaultConfig()
	return &HTTPService{
		svr: svr,
	}
}

// ListenAndServe is a function to start http service
func (h *HTTPService) ListenAndServe() error {
	return h.svr.ListenGracefully(":8080")
}
