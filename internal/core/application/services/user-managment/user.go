package usermanagment

import (
	"context"
	"time"

	"github.com/kaasikodes/assessmate_backend/internal/ports/outbound/email"
	jwtport "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/jwt"
	user_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/user"
)

type UserManagementService struct {
	userRepo user_repo.UserRepository
	jwt      *jwtport.JwtMaker
	//jwt
	//email
	emailClient email.EmailClient
}
type (
	LoginResponse struct {
		User         User
		AccessToken  string
		Institutions []Institution
	}
	User struct {
		Id        int
		Name      string
		Email     string
		Status    string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	Institution struct {
		Id          int
		Name        string
		Description string
		IsUserAdmin bool // is the user an admin
		IsActive    bool //represents whether user has access to institution could be blackisted or subscription has expired
	}
	ForgotPasswordResponse struct {
		Token string
		Email string
	}
	UserFilter struct {
		Status string
	}
	GetUsersResponse struct {
		Users []User
		Total int
	}
	AuthUserResponse struct {
		User         User
		Institutions []Institution
	}
)

// Constructor
func NewUserManagementService(repo user_repo.UserRepository) *UserManagementService {
	return &UserManagementService{
		userRepo: repo,
	}
}

// GetUserByEmail
func (u *UserManagementService) FindUserByEmail(ctx context.Context, email string) (*User, error) {

	return &User{}, nil
}

// login
func (u *UserManagementService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {

	return &LoginResponse{}, nil
}

// register
func (u *UserManagementService) Register(ctx context.Context, email, name, password string) (*User, error) {

	return &User{}, nil
}

// CreateToken
func (u *UserManagementService) CreateToken(ctx context.Context, userId int, tokenType string) (string, error) {

	return "", nil
}

// CreateToken
func (u *UserManagementService) GetToken(ctx context.Context, value, tokenType string) (string, error) {

	return "", nil
}

// DeleteToken
func (u *UserManagementService) DeleteToken(ctx context.Context, id int) error {

	return nil
}

// verifyAccount
func (u *UserManagementService) VerifyAccount(ctx context.Context, email, name, password string) (*LoginResponse, error) {

	return &LoginResponse{}, nil
}

// forgotPassword
func (u *UserManagementService) ForgotPassword(ctx context.Context, email string) (*ForgotPasswordResponse, error) {

	return &ForgotPasswordResponse{}, nil
}

// resetPassword
func (u *UserManagementService) ResetPassword(ctx context.Context, token, email, newPassword string) error {

	return nil
}

// getUsers
func (u *UserManagementService) GetUsers(ctx context.Context, filter *UserFilter) (*GetUsersResponse, error) {
	return &GetUsersResponse{}, nil

}

// getUserById
func (u *UserManagementService) FindUserById(ctx context.Context, id int) (*User, error) {

	return &User{}, nil
}

// logout
func (u *UserManagementService) Logout(ctx context.Context, email string) error {

	return nil
}

// getAuthUser
func (u *UserManagementService) GetAuthUser(ctx context.Context, email string) (*GetUsersResponse, error) {

	return &GetUsersResponse{}, nil
}
