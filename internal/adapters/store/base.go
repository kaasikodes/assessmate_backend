package store

import (
	"database/sql"

	// sub_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/subscription"
	"github.com/kaasikodes/assessmate_backend/internal/ports/outbound/logger"
	user_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/user"
)

type StoreCombinedRepository interface {
	user_repo.UserRepository
	// sub_repo.SubscriptionRepository
}
type MySqlRepo struct {
	db     *sql.DB
	logger logger.Logger
}

func NewUserRepository(db *sql.DB, logger logger.Logger) StoreCombinedRepository {
	return &MySqlRepo{db, logger}

}
