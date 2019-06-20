package server

import (
	"time"
)

type serverContainer struct {
	// servers is a map containing all the servers, the key will be the servers return from the GetUniqueIdentifier()
	servers map[string]Server
	// Err handles Any Errors that occure in the servercontainer will be sent to this Channel
	Err chan (error)
	// connectionTicker is a ticker that triggers the ConnectionTest function for all the servers
	connectionTicker *time.Ticker
	quit             chan (bool)
}

var defaultInterval = time.Duration(time.Second * 60)

// NewServerContainer will return a Correctly initialized serverContainer struct
func NewServerContainer() *serverContainer {
	return &serverContainer{
		servers: make(map[string]Server, 0),
		Err:     make(chan error),
		quit:    make(chan bool),
	}
}

// Add will take any Struct implementing the Server interface and add it among its servers
// Name is used for Map Lookup
func (sc *serverContainer) Add(s Server) error {
	/* Make sure we dont Overwrite any old server */
	if _, ok := sc.servers[s.GetUniqueIdentifier()]; ok {
		return ErrDupliaceServerName
	}
	sc.servers[s.GetUniqueIdentifier()] = s
	return nil
}

// Overwrite will add a Server at the UniqueIdentifier location even if there is already a Server
func (sc *serverContainer) Overwrite(s Server) {
	sc.servers[s.GetUniqueIdentifier()] = s
}

// Remove will delete the server from the ServerContainer
func (sc *serverContainer) Remove(s Server) {
	delete(sc.servers, s.GetUniqueIdentifier())
}

// SetConnectionInterval will change the duration between connection Checks on the servers
func (sc *serverContainer) SetConnectionInterval(sec int) {
	sc.connectionTicker = time.NewTicker(time.Duration(int64(sec)) * time.Second)
}

// StopManagingConnections will end the ManageConnections goroutine from running
func (sc *serverContainer) StopManagingConnections() {
	sc.connectionTicker.Stop()
	sc.quit <- true
}

// ManageConnections will itterate all current serveres and Check their connection state
// If they are unavailable it will try to Reconnect
// All errors will be sent to sc.Err
func (sc *serverContainer) ManageConnections() {
	go func() {
		if sc.connectionTicker == nil {
			/* Create a Default Intervalled Ticker */
			sc.connectionTicker = time.NewTicker(time.Duration(int64(defaultInterval)) * time.Second)
		}
		for {
			select {
			case <-sc.connectionTicker.C:
				for _, s := range sc.servers {
					go sc.manageConnection(s)
				}
			case <-sc.quit:
				return
			}
		}
	}()
}

// manageConnection will test the connection of a given Server
// Error will be reported on the Err Channel
func (sc *serverContainer) manageConnection(s Server) {
	sc.Err <- s.TestConnection()
}
