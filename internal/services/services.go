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
	DeleteUser(id int) error
	UpdateEmailUser(user *entity.User) (*entity.User, error)
	UpdatePasswordUser(user *entity.User) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	CheckUser(*entity.User) (*entity.User, error)
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
