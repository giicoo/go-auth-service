package user

import (
	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/repository"
)

type UserService struct {
	cfg *config.Config

	repo repository.Repo
}

func NewUserService(cfg *config.Config, repo repository.Repo) *UserService {
	return &UserService{
		cfg:  cfg,
		repo: repo,
	}
}
