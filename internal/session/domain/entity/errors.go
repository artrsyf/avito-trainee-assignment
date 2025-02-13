package entity

import "errors"

var (
	ErrNoSession        = errors.New("couldn't find session")
	ErrAlreadyCreated   = errors.New("session is already created")
	ErrWrongCredentials = errors.New("incorrect login or password")
)
