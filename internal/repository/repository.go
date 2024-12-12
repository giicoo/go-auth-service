package repository

import "github.com/giicoo/go-auth-service/internal/entity"

//go:generate mockgen -source=
type Repo interface {
	CreateUser(*entity.User) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	// UpdateUser(*entity.User) (*entity.User, error)
	// DeleteUser(id int) error
}
