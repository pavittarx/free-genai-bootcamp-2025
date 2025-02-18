package models

import "errors"

var (
	ErrInvalidID         = errors.New("invalid ID: must be a positive number")
	ErrInvalidTime       = errors.New("invalid time: time cannot be zero")
	ErrInvalidTimeRange  = errors.New("invalid time range: end time must be after start time")
	ErrInvalidScore      = errors.New("invalid score: score cannot be negative")
	ErrInvalidInput      = errors.New("invalid input: input cannot be empty")
)
