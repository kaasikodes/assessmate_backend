package store

import (
	"database/sql"

	sub_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/subscription"
	user_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/user"
)

type StoreCombinedRepository interface {
	user_repo.UserRepository
	sub_repo.SubscriptionRepository
}
type MySqlRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) StoreCombinedRepository {
	return &MySqlRepo{db}

}
