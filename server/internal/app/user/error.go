package user

import (
	"fmt"
)

type UnauthorizedError struct {
	Err error
}

func NewUnauthorizedError(err error) error {
	return &UnauthorizedError{Err: err}
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("%v : %v", e.Err, "already exists")
}
