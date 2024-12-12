package entity

import "errors"

type PublicError error
type PrivateError error

var (
	ErrInternal          PublicError = errors.New("internal error")
	ErrUserAlreadyExists PublicError = errors.New("user already exist")
)
