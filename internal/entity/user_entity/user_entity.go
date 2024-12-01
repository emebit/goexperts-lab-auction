package user_entity

import (
	"github.com/emebit/goexperts-lab-auction/internal/internal_error"
	"context"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
