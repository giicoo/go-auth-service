package sqlite

import "github.com/giicoo/go-auth-service/internal/config"

type Repo struct {
	cfg *config.Config
}

func NewRepo(cfg *config.Config) *Repo {
	return &Repo{
		cfg: cfg,
	}
}
