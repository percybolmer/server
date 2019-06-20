package server

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// basicServer is a Example server to match the interface, this struct should be reWritten to match your system type
type basicServer struct {
	id          string
	credentials *Credentials
	db          *gorm.DB
}

// NewBasicServer is a function used to initailze a new basicServer object, it will make sure
// that all neccesarry fields are initialized to avoid nil pointer exceptions
// id Should be Unique for each server
func NewBasicServer(id string) *basicServer {
	return &basicServer{
		id:          id,
		credentials: NewCredentials(),
	}
}

// Connect will start a connection against the Server
// Connect will be responsible to assign the needed values to the Server struct.
// Etc, If you run a REST server, add a Http.Client or so to the Server struct and make sure to assign it in the Connect upon success
// This BasicServer is a database server, so we will assign a gorm.DB object to basicserver.db
func (b *basicServer) Connect(c *Credentials) error {
	/* This example is a Database Application, so we need to validate that DB in Credentials is not nil */
	if c.DB == nil && b.credentials.DB == nil {
		return ErrMissingCredentials
	} else if c != nil {
		/* Assign New Credentials to server */
		b.credentials = c
	}
	dbinfo := b.credentials.DB
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbinfo.Host, dbinfo.Port, dbinfo.Username, dbinfo.Password, dbinfo.DatabaseName)
	db, err := gorm.Open("postgres", psqlInfo)
	db.LogMode(true)
	if err != nil {
		return ErrFailedToConnect

	}
	/* Assign the database into the struct */
	b.db = db
	return nil
}

// Reconnect is responsible to connect again, It can be a one call function that only calls Connect.
// However, Making sure the old connection is disconnected can be a good thing to do
func (b *basicServer) Reconnect() error {
	err := b.Disconnect()
	if err != nil && err != ErrCannotCloseNilDatabase {
		return err
	}
	return b.Connect(b.credentials)
}

// Disconnect is responble for closing the connection to the database
// Can be a good idea to make sure the object is not nil first.
func (b *basicServer) Disconnect() error {
	if b.db == nil {
		return ErrCannotCloseNilDatabase
	}
	return b.db.Close()
}

// TestConnection is a function to check that the Connection to the server is working properly
// This can be handled or just reported
// In this basicServer I have made the decission to call Reconnect and return THAT result if there is any failures
// If Connection is working properly, then just return nil
func (b *basicServer) TestConnection() error {
	if b.db == nil {
		return b.Reconnect()
	}
	/* Ask Database if connection is OK */
	err := b.db.DB().Ping()
	if err != nil {
		return b.Reconnect()
	}
	return nil
}

// Ping is used to just send a Ping to a DB, or a check response time on a HTTP Api etc,
func (b *basicServer) Ping() error {
	return b.db.DB().Ping()
}

// GetUniqueIdentifer returns the ID of the server
func (b *basicServer) GetUniqueIdentifier() string {
	return b.id
}

/** Add your own Functions down below */
