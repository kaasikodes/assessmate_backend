package user

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Id
type Id int

func (id Id) IsValid() bool {
	if id == 0 {
		return false

	}
	if id < 0 {
		return false
	}

	return true
}
func NewId(id int) (Id, error) {
	if id == 0 {
		return 0, errors.New("id cannot be zero")

	}
	if id < 0 {
		return 0, errors.New("id cannot be negative")
	}

	return Id(id), nil
}

// This is a convenient getter method.
func (i Id) Value() int {
	return int(i)
}

// UserStatus
type UserStatus string

var (
	InActive UserStatus = "inactive"
	Active   UserStatus = "active"
)

func NewUserStatus(val string) (UserStatus, error) {
	if isValidUserStatus(val) {
		return UserStatus(val), nil

	}
	return "", errors.New("the user staus is not recognized")

}
func (st UserStatus) IsValid() bool {
	return isValidUserStatus(string(st))
}
func (st UserStatus) String() string {
	return string(st)
}

// IsValid checks if the UserStatus is one of the predefined valid types.
func isValidUserStatus(val string) bool {
	switch UserStatus(val) {
	case Active, InActive:
		return true
	default:
		return false
	}
}

type DateTime = time.Time

func NewDateTime(input string) (*DateTime, error) {
	if input == "" {
		return nil, nil // empty string means optional/no filter
	}

	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}
	return &t, nil
}

type Name string

func (n Name) String() string {
	return string(n)

}
func NewName(name string) (Name, error) {
	// cannot be empty
	name = strings.TrimSpace(name)

	if name == "" {
		return "", errors.New("name cannot be empty")
	}
	minLength := 3
	maxLength := 100

	// must be at least 3 characters long
	if len(name) < minLength {
		return "", (fmt.Errorf("name must be at least %d charaters long", minLength))

	}
	// must be at least max characters long
	if len(name) > maxLength {
		return "", (fmt.Errorf("name must not exceed %d charaters", maxLength))

	}

	return Name(name), nil

}
func (n Name) IsEmpty() bool {
	return strings.TrimSpace(string(n)) == ""

}

// TokenType
type TokenType string

var (
	Verification  TokenType = "verification"
	ResetPassword TokenType = "reset-password"
	RefreshToken  TokenType = "refresh-token"
)

func NewTokenType(val string) (TokenType, error) {
	if isValidTokenType(val) {
		return TokenType(val), nil

	}
	return "", errors.New("the token type is not recognized")

}
func (st TokenType) IsValid() bool {
	return isValidTokenType(string(st))
}
func (st TokenType) String() string {
	return string(st)
}

// IsValid checks if the TokenType is one of the predefined valid types.
func isValidTokenType(val string) bool {
	switch TokenType(val) {
	case RefreshToken, ResetPassword, Verification:
		return true
	default:
		return false
	}
}

type TokenValue string

func (n TokenValue) String() string {
	return string(n)

}
func NewTokenValue(val string) (TokenValue, error) {
	// cannot be empty
	val = strings.TrimSpace(val)

	if val == "" {
		return "", errors.New("token value cannot be empty")
	}
	minLength := 3
	maxLength := 20

	// must be at least 3 characters long
	if len(val) < minLength {
		return "", (fmt.Errorf("name must be at least %d charaters long", minLength))

	}
	// must be at least max characters long
	if len(val) > maxLength {
		return "", (fmt.Errorf("name must not exceed %d charaters", maxLength))

	}

	return TokenValue(val), nil

}

func (n TokenValue) IsEmpty() bool {
	return strings.TrimSpace(string(n)) == ""

}
