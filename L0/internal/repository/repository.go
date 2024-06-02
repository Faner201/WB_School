package repository

import (
	"L0/internal/util"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db     *sqlx.DB
	config *util.Config
	cache  *Cache
}

func NewRepository(config *util.Config, cache *Cache) *Repository {
	return &Repository{
		config: config,
		cache:  cache,
	}
}
