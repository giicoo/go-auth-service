package services

import (
	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/repository"
	"github.com/giicoo/go-auth-service/internal/services/user"
)

type UserService interface {
	CreateUser(*entity.User) (*entity.User, error)
}

type Services struct {
	cfg *config.Config

	UserService UserService
}

func NewServices(cfg *config.Config, repo repository.Repo) *Services {
	return &Services{
		UserService: user.NewUserService(cfg, repo),
	}
}
