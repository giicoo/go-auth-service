package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/giicoo/go-auth-service/internal/config"
	"github.com/giicoo/go-auth-service/internal/entity"
)

type Repo struct {
	cfg *config.Config

	db *sql.DB
}

func NewRepo(cfg *config.Config, db *sql.DB) *Repo {
	return &Repo{
		cfg: cfg,

		db: db,
	}
}

func (r *Repo) InitRepo() error {
	file, err := os.ReadFile(r.cfg.DB.PathToSQL + "create_table.sql")
	if err != nil {
		return fmt.Errorf("read create sql file: %w", err)
	}
	stmt := string(file)
	_, err = r.db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	err = r.db.Ping()
	if err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	return nil
}

func (r *Repo) CreateUser(user *entity.User) (*entity.User, error) {
	stmt := "INSERT INTO users (email, hash_password) VALUES (?, ?)"

	if _, err := r.db.Exec(stmt, user.Email, user.Password); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}

	userDB, err := r.GetUserByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("get user by %s: %w", user.Email, err)
	}

	return userDB, nil
}

func (r *Repo) GetUserByEmail(email string) (*entity.User, error) {
	stmt := "SELECT * FROM users WHERE email=?"

	userDB := new(entity.User)

	row := r.db.QueryRow(stmt, email)
	err := row.Scan(&userDB.ID, &userDB.Email, &userDB.Password)
	if err != nil {
		return nil, fmt.Errorf("row scan: %w", err)
	}

	return userDB, nil
}
