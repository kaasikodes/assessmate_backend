package usermanagment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/user"
	email_client "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/email"
	jwtport "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/jwt"
	user_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/user"
)

type UserManagementService struct {
	userRepo user_repo.UserRepository
	jwt      jwtport.JwtMaker
	//jwt
	//email
	emailClient email_client.EmailClient
}
type (
	LoginResponse struct {
		User         User
		AccessToken  string
		Institutions []Institution
	}
	User struct {
		Id         int
		Name       string
		Email      string
		Status     string
		IsVerified bool
		CreatedAt  time.Time
		UpdatedAt  time.Time
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
		Status *string
	}
	GetUsersResponse struct {
		Users []User
		Total int
	}
	AuthUserResponse struct {
		User         User
		Institutions []Institution
	}
	TokenResponse struct {
		Id        int
		Value     string
		TokenType string
	}
)

// Constructor
func NewUserManagementService(repo user_repo.UserRepository, jwt jwtport.JwtMaker, emailClient email_client.EmailClient) *UserManagementService {
	return &UserManagementService{
		userRepo:    repo,
		jwt:         jwt,
		emailClient: emailClient,
	}
}

// GetUserByEmail
func (u *UserManagementService) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	parsedEmail, err := user.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error parsing email: %w", err)
	}
	domainUser, err := u.userRepo.GetUserByEmail(ctx, parsedEmail)
	if err != nil {
		return nil, err
	}

	return mapToServiceUser(domainUser), nil
}

func (u *UserManagementService) FindUserById(ctx context.Context, id int) (*User, error) {
	parsedId, err := user.NewId(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %w", err)
	}
	domainUser, err := u.userRepo.GetUserById(ctx, parsedId)
	if err != nil {
		return nil, err
	}

	return mapToServiceUser(domainUser), nil
}

func (u *UserManagementService) createAccessToken(userId, userEmail string) (string, error) {
	tokenDuration := time.Hour * 24 * 3 //3 days
	token, err := u.jwt.CreateToken(userId, userEmail, tokenDuration)
	if err != nil {
		return "", err
	}
	return token, nil
}

// login
func (u *UserManagementService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	parsedEmail, err := user.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error parsing email: %w", err)
	}

	domainUser, err := u.userRepo.GetUserByEmail(ctx, parsedEmail)
	if err != nil {
		return nil, err
	}
	passwordsMatch := domainUser.ComparePassword(password)
	if !passwordsMatch {
		return nil, errors.New("invalid credentials")
	}
	token, err := u.createAccessToken(domainUser.GetId().String(), domainUser.GetEmail().String())
	if err != nil {
		return nil, fmt.Errorf("error creating access token: %w", err)
	}
	return &LoginResponse{
		User:         *mapToServiceUser(domainUser),
		AccessToken:  token,
		Institutions: []Institution{
			// TODO: Get Instititutions user belongs to if any and populate here, might have to be a bounded context or something or take in institute_repo and update
		},
	}, nil
}

// register
func (u *UserManagementService) Register(ctx context.Context, email, name, password string) (*User, error) {
	parsedName, err := user.NewName(name)
	if err != nil {
		return nil, fmt.Errorf("error parsing name: %w", err)
	}
	parsedEmail, err := user.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error parsing email: %w", err)
	}
	domainUser, err := user.NewUser(parsedName, parsedEmail)
	if err != nil {
		return nil, fmt.Errorf("error parsing user: %w", err)
	}
	err = domainUser.SetPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error setting user password: %w", err)
	}

	createdUser, err := u.userRepo.CreateUser(ctx, domainUser)
	if err != nil {
		return nil, err
	}
	// Create verification token and send to user via mail
	tokenVal, err := user.NewTokenValue("verify_" + time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error parsing token value: %w", err)
	}

	verifyToken, err := user.NewToken(tokenVal, user.ResetPassword, domainUser.GetId())
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}
	verifyToken, err = u.userRepo.CreateToken(ctx, verifyToken.Value(), user.ResetPassword, domainUser.GetId())
	if err != nil {
		return nil, fmt.Errorf("error saving token: %w", err)
	}

	go func() {
		// Send email via email client
		err = u.emailClient.Send(ctx, &email_client.Notification{

			Email:   createdUser.GetEmail().String(),
			Title:   "Account Verification",
			Content: fmt.Sprintf("This is your verification token: %s", verifyToken.Value().String()),
		})

	}()

	return mapToServiceUser(createdUser), nil
}

// CreateToken
func (u *UserManagementService) CreateToken(ctx context.Context, userId int, tokenType string) (string, error) {
	tokenVal, err := user.NewTokenValue("tok_" + time.Now().Format(time.RFC3339))
	if err != nil {
		return "", fmt.Errorf("error constructing token value: %w", err)
	}
	parsedTokenType, err := user.NewTokenType(tokenType)
	if err != nil {
		return "", fmt.Errorf("error parsing token type: %w", err)
	}
	parsedUserId, err := user.NewId(userId)
	if err != nil {
		return "", fmt.Errorf("error parsing userId: %w", err)
	}

	_, err = u.userRepo.CreateToken(ctx, tokenVal, parsedTokenType, parsedUserId)
	if err != nil {
		return "", fmt.Errorf("error storing token: %w", err)
	}

	return tokenVal.String(), nil
}

// GetToken
func (u *UserManagementService) GetToken(ctx context.Context, value, tokenType string, userId int) (*TokenResponse, error) {
	tokenVal, err := user.NewTokenValue(value)
	if err != nil {
		return nil, fmt.Errorf("error constructing token value: %w", err)
	}
	parsedUserId, err := user.NewId(userId)
	if err != nil {
		return nil, fmt.Errorf("error parsing userId: %w", err)
	}
	token, err := u.userRepo.GetToken(ctx, parsedUserId, tokenVal)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		Id:        token.Id().Value(),
		Value:     token.Value().String(),
		TokenType: token.Type().String(),
	}, nil
}

// DeleteToken
func (u *UserManagementService) DeleteToken(ctx context.Context, id int) error {
	tokenId, err := user.NewId(id)
	if err != nil {
		return fmt.Errorf("error parsing tokenId: %w", err)
	}

	return u.userRepo.DeleteToken(ctx, tokenId)
}

// verifyAccount
func (u *UserManagementService) VerifyUser(ctx context.Context, email, _token string) (*LoginResponse, error) {
	tokenVal, err := user.NewTokenValue(_token)
	if err != nil {
		return nil, fmt.Errorf("error constructing token value: %w", err)
	}
	parsedEmail, err := user.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error parsing email: %w", err)
	}
	domainUser, err := u.userRepo.GetUserByEmail(ctx, parsedEmail)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	token, err := u.userRepo.GetToken(ctx, domainUser.GetId(), tokenVal)
	if err != nil {
		return nil, fmt.Errorf("error retrieving token: %w", err)
	}

	domainUser.SetVerifiedAt(time.Now())

	domainUser, err = u.userRepo.VerifyUser(ctx, domainUser)
	if err != nil {
		return nil, fmt.Errorf("error verifying user: %w", err)
	}
	// delete the token
	err = u.userRepo.DeleteToken(ctx, token.Id())
	if err != nil {
		return nil, fmt.Errorf("error deleting token: %w", err)
	}

	accessToken, err := u.createAccessToken(domainUser.GetId().String(), domainUser.GetEmail().String())
	if err != nil {
		return nil, fmt.Errorf("error creating access token: %w", err)
	}

	return &LoginResponse{
		User:        *mapToServiceUser(domainUser),
		AccessToken: accessToken,
	}, nil
}

// forgotPassword
func (u *UserManagementService) ForgotPassword(ctx context.Context, email string) (*ForgotPasswordResponse, error) {

	domainUser, err := u.userRepo.GetUserByEmail(ctx, user.Email(email))
	if err != nil {
		return nil, err
	}
	tokenVal, err := user.NewTokenValue("reset_" + time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("error parsing token value: %w", err)
	}

	resetToken, err := user.NewToken(tokenVal, user.ResetPassword, domainUser.GetId())
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	resetToken, err = u.userRepo.CreateToken(ctx, resetToken.Value(), user.ResetPassword, domainUser.GetId())
	if err != nil {
		return nil, fmt.Errorf("error saving token: %w", err)
	}

	go func() {
		// Send email via email client
		err = u.emailClient.Send(ctx, &email_client.Notification{
			Title: "Password Reset",
			Email: domainUser.GetEmail().String(),

			Content: fmt.Sprintf("This is your password reset token: %s", resetToken.Value().String()),
		})

	}()

	return &ForgotPasswordResponse{
		Token: resetToken.Value().String(),
		Email: domainUser.GetEmail().String(),
	}, nil
}

// resetPassword
func (u *UserManagementService) ResetPassword(ctx context.Context, _token, email, newPassword string, userId int) error {
	tokenVal, err := user.NewTokenValue(_token)
	if err != nil {
		return fmt.Errorf("error constructing token value: %w", err)
	}
	parsedUserId, err := user.NewId(userId)
	if err != nil {
		return fmt.Errorf("error parsing userId: %w", err)
	}
	token, err := u.userRepo.GetToken(ctx, parsedUserId, tokenVal)
	if err != nil {
		return fmt.Errorf("error retrieving token: %w", err)
	}

	domainUser, err := u.userRepo.GetUserById(ctx, parsedUserId)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	domainUser.SetPassword(newPassword)

	_, err = u.userRepo.UpdateUserPassword(ctx, domainUser)
	if err != nil {
		return fmt.Errorf("error updating user password: %w", err)
	}

	// delete the token
	err = u.userRepo.DeleteToken(ctx, token.Id())
	if err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}

	return nil
}

// getUsers
func (u *UserManagementService) GetUsers(ctx context.Context, filter *UserFilter) (*GetUsersResponse, error) {
	var parsedFilter user.UserFilter
	if filter != nil && filter.Status != nil {
		status, err := user.NewUserStatus(*filter.Status)
		if err != nil {
			return nil, fmt.Errorf("error parsing user status: %w", err)
		}
		parsedStatus := status
		parsedFilter.Status = &parsedStatus

	}
	data, total, err := u.userRepo.GetUsers(ctx, &parsedFilter)
	if err != nil {
		return nil, fmt.Errorf("error retrieving users from store: %w", err)
	}

	users := make([]User, len(data))
	for i, u := range data {

		users[i] = *mapToServiceUser(&u)

	}

	return &GetUsersResponse{
		Users: users,
		Total: total,
	}, nil

}

// logout
func (u *UserManagementService) Logout(ctx context.Context, email string) error {
	//TODO: create an audit repo that will record the logout event

	return nil
}

// getAuthUser
func (u *UserManagementService) GetAuthUser(ctx context.Context, email string) (*AuthUserResponse, error) {
	parsedEmail, err := user.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error parsing email: %w", err)
	}
	domainUser, err := u.userRepo.GetUserByEmail(ctx, parsedEmail)
	if err != nil {
		return nil, err
	}

	return &AuthUserResponse{
		User: *mapToServiceUser(domainUser),
	}, nil
}

// Helpers
func mapToServiceUser(domainUser *user.User) *User {
	return &User{
		Id:         domainUser.GetId().Value(),
		Name:       domainUser.GetName().String(),
		Email:      domainUser.GetEmail().String(),
		Status:     domainUser.GetStatus().String(),
		CreatedAt:  domainUser.GetCreatedAt(),
		UpdatedAt:  domainUser.GetUpdatedAt(),
		IsVerified: domainUser.IsVerified(),
	}
}
