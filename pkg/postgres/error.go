package postgres

import "errors"

var (
	ErrBuildingQuery = errors.New("building query")
	ErrReadingRow    = errors.New("reading row")
	ErrNoUpdates     = errors.New("no fields to update")
)
