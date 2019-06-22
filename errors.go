package server

import "errors"

var (
	ErrFailedToConnect error = errors.New("Failed to connect to server")

	ErrCannotCloseNilDatabase error = errors.New("The database object is nil, cannot call close")

	/* ServerContainer errors */
	ErrDupliaceServerName error = errors.New("There is already a Server with this name, to overwrite please use Overwrite function")
)
