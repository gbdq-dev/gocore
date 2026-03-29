// Package http provides a generic HTTP server implementation that works with any
// router/framework that implements the http.Handler interface.
package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// GenericServer represents an HTTP server that can work with any router compatible with standard http package.
type GenericServer struct {
	Port   string       // The port on which the server will run.
	Server *http.Server // The underlying HTTP server.
	log    *slog.Logger // Logger for server operations.

	InitRoutes func(http.Handler) // A function to initialize routes for the server.
}

// NewGenericServer creates a new generic HTTP server instance.
//
// Arguments:
//
//	port - The port on which the server will run.
//	server - The underlying HTTP server instance.
//	logger - Logger for server operations.
//	initRoutes - A function to initialize routes for the server. This function should take an http.Handler as an argument and set up the routes for the server. If nil, no routes will be initialized.
//
// Returns:
//
//	A pointer to a new GenericServer instance.
func NewGenericServer(port string, server *http.Server, logger *slog.Logger, initRoutes func(http.Handler)) *GenericServer {
	if initRoutes != nil {
		initRoutes(server.Handler) // Initialize routes if a function is provided.
	}

	return &GenericServer{
		Port:   port,
		Server: server,
		log:    logger,
	}
}

// Start starts the HTTP server on the specified port.
//
// Returns:
//
//	An error if the server fails to start.
func (s *GenericServer) Start() error {
	s.log.Info("HTTP server running on port " + s.Port)
	return s.Server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server.
//
// This method uses the standard http.Server's Shutdown method to ensure all ongoing
// requests are allowed to complete or time out before the server shuts down.
//
// Returns:
//
//	An error if the server fails to shut down gracefully.
func (s *GenericServer) Stop() error {
	s.log.Info("HTTP server stopping gracefully")

	// Create a context with timeout for the shutdown operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	return s.Server.Shutdown(ctx)
}
