package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/repository"
	"github.com/giicoo/go-auth-service/pkg/apiError"
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
		return nil, fmt.Errorf("user`s email %s: %w", userYet.Email, apiError.ErrUserAlreadyExists)
	}

	hashPassword, err := hashTools.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashPassword

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	userDB, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("get user by %s email: %w", user.Email, err)
	}
	return userDB, nil
}
func (s *UserService) CheckUser(user *entity.User) (*entity.User, error) {
	userDB, err := s.repo.GetUserByEmail(user.Email)
	if errors.Is(err, sql.ErrNoRows) && err != nil {
		return nil, fmt.Errorf("user email %s: %w", user.Email, apiError.ErrUserNotExists)
	}
	if err != nil {
		return nil, fmt.Errorf("user check %s: %w", user.Email, err)
	}

	if !hashTools.CheckPasswordHash(user.Password, userDB.Password) {
		return nil, fmt.Errorf("check user %s: %w", user.Email, apiError.ErrWrongPassword)
	}

	return userDB, nil
}
func (s *UserService) DeleteUser(id int) error {
	userYet, err := s.repo.GetUserByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get user by %d id: %w", id, apiError.ErrUserNotExists)
	}
	if err != nil {
		return fmt.Errorf("get user by %d id: %w", id, err)
	}

	if err := s.repo.DeleteUser(userYet.ID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (s *UserService) UpdateEmailUser(user *entity.User) (*entity.User, error) {
	_, err := s.repo.GetUserByID(user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get user by %d id: %w", user.ID, apiError.ErrUserNotExists)
	}
	if err != nil {
		return nil, fmt.Errorf("get user by %d id: %w", user.ID, err)
	}

	// check exist email
	userYet, err := s.repo.GetUserByEmail(user.Email)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return nil, fmt.Errorf("check exist user: %w", err)
	}
	if userYet != nil {
		return nil, fmt.Errorf("user`s email %s: %w", userYet.Email, apiError.ErrUserAlreadyExists)
	}

	if err := s.repo.UpdateEmailUser(user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	userDB, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get user by %s email: %w", user.Email, err)
	}
	return userDB, nil
}

func (s *UserService) UpdatePasswordUser(user *entity.User) (*entity.User, error) {
	_, err := s.repo.GetUserByID(user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get user by %d id: %w", user.ID, apiError.ErrUserNotExists)
	}
	if err != nil {
		return nil, fmt.Errorf("get user by %d id: %w", user.ID, err)
	}

	hashPassword, err := hashTools.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashPassword

	if err := s.repo.UpdatePasswordUser(user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	userDB, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get user by %s email: %w", user.Email, err)
	}
	return userDB, nil
}

// func (s *UserService) GetUserByEmail(user *entity.User) (*entity.User, error) {
// 	userDB, err := s.repo.GetUserByEmail(user.Email)
// 	if errors.Is(err, sql.ErrNoRows) {
// 		return nil, fmt.Errorf("get user by %s email: %w", user.Email, apiError.ErrUserNotExists)
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("get user by %s email: %w", user.Email, err)
// 	}
// 	return userDB, nil
// }

func (s *UserService) GetUserByID(id int) (*entity.User, error) {
	userDB, err := s.repo.GetUserByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get user by %d id: %w", id, apiError.ErrUserNotExists)
	}
	if err != nil {
		return nil, fmt.Errorf("get user by %d id: %w", id, err)
	}
	return userDB, nil
}
