package storage

import "errors"

var (
	ErrAPODNotFound = errors.New("APOD not found")
	ErrAPODExists   = errors.New("APOD exists")
)
