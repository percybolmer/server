package server

import "errors"

var (
	ErrFailedToConnect error = errors.New("Failed to connect to server")
	/* Credentials realted errors */
	ErrMissingCredentials           error = errors.New("Missing the needed credentials for this connection type")
	ErrFailedToUnmarshalCredentials error = errors.New("Failed to unmarshal provided JSON data")
	ErrNeedUserAndPassword          error = errors.New("Need both username or password")

	ErrInvalidHost   error = errors.New("The host provided is invalid")
	ErrSshInvalidKey error = errors.New("The provided key file seems to be invalid")

	ErrInvalidDatabaseName    error = errors.New("The database name is invalid")
	ErrInvalidDatabasePort    error = errors.New("Database port is invalids")
	ErrCannotCloseNilDatabase error = errors.New("The database object is nil, cannot call close")

	/* ServerContainer errors */
	ErrDupliaceServerName error = errors.New("There is already a Server with this name, to overwrite please use Overwrite function")
)
