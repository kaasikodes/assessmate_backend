package user

import (
	"context"
	"errors"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/user"
)

var ErrUserNotFound = errors.New("user not found")

// ports should conform to language of core(in this case the domain and not application, as application is a bridge for adapter to domain(business) logic)
type UserRepository interface {
	CreateUser(ctx context.Context, payload *user.User) (*user.User, error)
	ChangeUserStatus(ctx context.Context, userId user.User, status *user.UserStatus) (*user.User, error)
}
