package user

import (
	"errors"
	"time"
)

type Token struct {
	id        Id
	value     TokenValue
	tokenType TokenType
	userId    Id
	createdAt DateTime
	updatedAt DateTime
	expiresAt *DateTime
}

// NewToken creates a new token instance.
func NewToken(value TokenValue, tokenType TokenType, userId Id) (*Token, error) {

	if value.IsEmpty() {
		return nil, errors.New("token value cannot be empty")
	}
	if !userId.IsValid() {
		return nil, errors.New("user id cannot be 0 or negative")
	}
	if !tokenType.IsValid() {
		return nil, errors.New("invalid token type")
	}

	now := DateTime(time.Now().UTC())

	return &Token{
		id:        Id(0),
		value:     value,
		tokenType: tokenType,
		userId:    userId,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// SetId sets the token ID if not already set.
func (t *Token) SetId(id Id) {

	t.id = id
	t.touch()
}

// SetUserId
func (t *Token) SetUserId(userId Id) {
	t.userId = userId
	t.touch()
}

// SetValue sets a new token value and updates the updatedAt timestamp.
func (t *Token) SetValue(value TokenValue) {
	t.value = value
	t.touch()
}

// SetValue sets a new token value and updates the updatedAt timestamp.
func (t *Token) SetExpiresAt(ti time.Time) {
	d := DateTime(ti)
	t.expiresAt = &d
	t.touch()
}

// SetCreatedAt manually updates the timestamp.
func (t *Token) SetType(ti string) {
	t.tokenType = TokenType(ti)
}

// SetCreatedAt manually updates the timestamp.
func (t *Token) SetCreatedAt(ti time.Time) {
	t.createdAt = DateTime(ti)
}

// SetUpdatedAt manually updates the timestamp.
func (t *Token) SetUpdatedAt(ti time.Time) {
	t.updatedAt = DateTime(ti)
}

// Getters
func (t *Token) Id() Id {
	return t.id
}

func (t *Token) Value() TokenValue {
	return t.value
}

func (t *Token) Type() TokenType {
	return t.tokenType
}

func (t *Token) UserId() Id {
	return t.userId
}

func (t *Token) CreatedAt() DateTime {
	return t.createdAt
}

func (t *Token) UpdatedAt() DateTime {
	return t.updatedAt
}
func (t *Token) ExpiresAt() *DateTime {
	return t.expiresAt
}
func (t *Token) HasExpired() bool {
	return t.ExpiresAt().Before(time.Now())
}

// touch updates the updatedAt timestamp.
func (u *Token) touch() {
	u.updatedAt = DateTime(time.Now().UTC())
}
