package entity

import "errors"

var (
	ErrNotEnoughBalance  = errors.New("not enough balance")
	ErrNotExistedProduct = errors.New("item is missing from the store")
)
