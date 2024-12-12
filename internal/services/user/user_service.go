package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/repository"
	hashTools "github.com/giicoo/go-auth-service/pkg/hash"
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

func (s *UserService) CreateUser(user *entity.User) (*entity.User, error) {
	userYet, err := s.repo.GetUserByEmail(user.Email)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return nil, fmt.Errorf("check exist user: %w", err)
	}
	if userYet != nil {
		return nil, fmt.Errorf("user`s email %s: %w", userYet.Email, entity.ErrUserAlreadyExists)
	}

	hashPassword, err := hashTools.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashPassword

	userDB, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return userDB, nil
}
