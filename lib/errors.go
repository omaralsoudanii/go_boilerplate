package lib

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrUnauthorized        = errors.New("You are unauthorized to use this")
	ErrNotFound            = errors.New("Your requested entity is not found")
	ErrConflict            = errors.New("Your requested entity already exist")
	ErrBadParamInput       = errors.New("Given Param is not valid")
	ErrInvalidTkn          = errors.New("Malformed authentication token")
	ErrInvalidRefreshTkn   = errors.New("Malformed authentication refresh token")
	ErrBadGrantType        = errors.New("Unsupported grant type")
)
