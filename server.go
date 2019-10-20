// Package server is used to ease the adding / creation of a new Server type so that they can be Combined easily
// @Author Perbol
// @Version 1.0
package server

import (
	"github.com/percybolmer/credentials"
)

// Server is the interface used to handle Servers
type Server interface {
	Connect(c *credentials.Credentials) error
	Reconnect() error
	Disconnect() error
	TestConnection() error
	Ping() error
	// GetUniqueIdentifier is responsible to return a Unique identifier for each Server object
	GetUniqueIdentifier() string
}
