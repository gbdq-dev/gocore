// Package http provides an HTTP server implementation using the Fiber framework.
// It allows for easy initialization, route setup, and server management.
package http

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Server represents an HTTP server using the Fiber framework.
type Server struct {
	Port      string               // The port on which the server will run.
	App       *fiber.App           // The Fiber app instance.
	InitRoute func(app *fiber.App) // A function to initialize routes for the Fiber app.

	log *slog.Logger
}

// NewFiber creates a new HTTP server instance using the Fiber framework.
//
// Arguments:
//
//	port - The port on which the server will run.
//	app - The Fiber app instance to be used.
//	initRoute - A function for initializing routes in the Fiber app (can be nil).
//
// Returns:
//
//	A pointer to a new Server instance.
func NewFiber(port string, app *fiber.App, initRoute func(app *fiber.App), logger *slog.Logger) *Server {
	if initRoute != nil {
		initRoute(app) // Initialize routes if a function is provided.
	}
	return &Server{
		Port:      port,
		App:       app,
		InitRoute: initRoute,
		log:       logger,
	}
}

// Start starts the HTTP server on the specified port.
//
// Returns:
//
//	An error if the server fails to start.
func (s *Server) Start() error {
	s.log.Info("HTTP server running on port " + s.Port)
	return s.App.Listen(":" + s.Port)
}

// Stop gracefully shuts down the HTTP server.
//
// This implementation uses Fiber's Shutdown method to ensure all ongoing
// requests are allowed to complete or time out before the server shuts down.
//
// Returns:
//
//	An error if the server fails to shut down gracefully.
func (s *Server) Stop() error {
	s.log.Info("HTTP server stopping gracefully")

	// Create a context with timeout for the shutdown operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the Fiber app gracefully
	return s.App.ShutdownWithContext(ctx)
}
