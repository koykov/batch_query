package sql

import "errors"

var (
	ErrNoDB       = errors.New("no database provided")
	ErrNoQuery    = errors.New("no query provided")
	ErrNoQueryFmt = errors.New("no query formatter provided")
	ErrNoRecScnr  = errors.New("no record scanner provided")

	ErrUnknownPlaceholderType = errors.New("unknown placeholder type provided")
)
