package connection

import "errors"

// Errors
var (
	ErrConnectionNotFound    = errors.New("connection not found")
	ErrConnectionAlreadyOpen = errors.New("connection is already open")
	ErrConnectionClosed      = errors.New("connection is not open")
	ErrInvalidParameters     = errors.New("invalid connection parameters")
)
