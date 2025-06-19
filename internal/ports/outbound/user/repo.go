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
	UpdateUserPassword(ctx context.Context, payload *user.User) (*user.User, error)
	VerifyUser(ctx context.Context, payload *user.User) (*user.User, error)
	ChangeUserStatus(ctx context.Context, userId user.Id, status user.UserStatus) (*user.User, error)
	CreateToken(ctx context.Context, value user.TokenValue, tokenType user.TokenType, userId user.Id, expiresAt user.DateTime) (*user.Token, error)
	DeleteToken(ctx context.Context, id user.Id) error
	GetToken(ctx context.Context, userId user.Id, value user.TokenValue) (*user.Token, error)
	GetUserById(ctx context.Context, userId user.Id) (*user.User, error)
	GetUserByEmail(ctx context.Context, user user.Email) (*user.User, error)
	GetUsers(ctx context.Context, filter *user.UserFilter) ([]user.User, int, error)
}
