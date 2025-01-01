package services

import (
	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/repository"
	"github.com/giicoo/go-auth-service/internal/services/session"
	"github.com/giicoo/go-auth-service/internal/services/user"
)

//go:generate mockgen -source=services.go -destination=mocks/mock.go
type UserService interface {
	CreateUser(*entity.User) (*entity.User, error)
	DeleteUser(*entity.User) error
	UpdateUser(*entity.User) (*entity.User, error)
	GetUserByEmail(*entity.User) (*entity.User, error)
	GetUserByID(*entity.User) (*entity.User, error)
}

type Services struct {
	cfg *config.Config

	UserService    UserService
	SessionService *session.SessionService
}

func NewServices(cfg *config.Config, repo repository.Repo) *Services {
	return &Services{
		UserService:    user.NewUserService(cfg, repo),
		SessionService: session.NewSessionManager(),
	}
}
