package postgres

import (
	"github.com/CrabRus/LiveStats/internal/db"
	"github.com/CrabRus/LiveStats/internal/domain"
)

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(database *db.DB) domain.UserRepository {
	return &UserRepository{db: database}
}
