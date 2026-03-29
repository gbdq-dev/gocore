// Package servers provides utilities to manage multiple server instances,
// including starting and stopping them in a concurrent and orderly manner.
package servers

import (
	"log/slog"
	"sync"
)

// Server defines the interface that any server must implement
// to be managed by the ServerManager.
type Server interface {
	Start() error // Starts the server. Returns an error if the server fails to start.
	Stop() error  // Stops the server. Returns an error if the server fails to stop.
}

// ServerManager manages a collection of servers, providing functionality
// to start and stop all servers concurrently.
type ServerManager struct {
	servers []Server     // List of servers to be managed.
	log     *slog.Logger // Logger for logging server events.
}

// NewServerManager creates a new instance of ServerManager.
//
// Returns:
//
//	A pointer to a newly initialized ServerManager.
func NewServerManager(logger *slog.Logger) *ServerManager {
	return &ServerManager{
		log: logger,
	}
}

// AddServer adds a new server to the ServerManager's list of servers.
//
// Arguments:
//
//	server - An instance of a Server to be added to the manager.
func (sm *ServerManager) AddServer(server Server) {
	sm.servers = append(sm.servers, server)
}

// StartAll starts all servers managed by the ServerManager concurrently.
// Any errors encountered while starting a server are logged.
func (sm *ServerManager) StartAll() {
	var wg sync.WaitGroup
	for _, server := range sm.servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			if err := s.Start(); err != nil {
				sm.log.Error("Failed to start server", "error", err)
			}
		}(server)
	}
	wg.Wait() // Wait for all servers to start.
}

// StopAll stops all servers managed by the ServerManager sequentially.
// Any errors encountered while stopping a server are logged.
func (sm *ServerManager) StopAll() error {
	for _, server := range sm.servers {
		err := server.Stop()
		if err != nil {
			return err
		}
	}

	return nil
}
