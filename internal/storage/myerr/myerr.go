package myerr

import "github.com/pkg/errors"

// Ошибки базы данных
var (
	FailedToCreateUserErr = errors.New("failed to create the user")
	UserNotFoundErr       = errors.New("user not found")
	AlreadyExistsErr      = errors.New("user already exists")
	NotFoundErr           = errors.New("not found")
)

var (
	ErrInvalidToken = errors.New("invalid token")
)
