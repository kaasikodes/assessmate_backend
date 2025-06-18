package jwtport

import (
	"errors"
	"net/http"
	"time"
)

type CustomClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email,omitempty"` // Optional field
}
type JwtMaker interface {
	// CreateToken generates a JWT signed with HS256
	CreateToken(userID, userEmail string, duration time.Duration) (string, error)

	// VerifyToken parses and validates the JWT token
	VerifyToken(tokenStr string) (*CustomClaims, error)

	// ExtractToken extracts token from Authorization header
	ExtractToken(r *http.Request) (string, error)

	// ExtractAndVerifyToken extracts the token from the request and verifies it
	ExtractAndVerifyToken(r *http.Request) (*CustomClaims, error)
}

var (
	ErrExpiredToken = errors.New("expired token")
	ErrInvalidToken = errors.New("invalid token")
	ErrNoAuthHeader = errors.New("authorization header is missing")
	ErrWrongFormat  = errors.New("authorization header format must be Bearer {token}")
)
