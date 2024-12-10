package services

import (
	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/repository"
	"github.com/giicoo/go-auth-service/internal/services/user"
)

type UserService interface {
	// CreateUser(entity.User) (entity.User, error)
	// UpdateUser(entity.User) (entity.User, error)
	// DeleteUser(id int) error
}

type Services struct {
	cfg *config.Config

	userService UserService
}

func NewServices(cfg *config.Config, repo repository.Repo) *Services {
	return &Services{
		userService: user.NewUserService(cfg, repo),
	}
}
