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
}

// NewToken creates a new token instance.
func NewToken(value TokenValue, tokenType TokenType, userId Id) (*Token, error) {

	if value.IsEmpty() {
		return nil, errors.New("token value cannot be empty")
	}
	if userId.IsValid() {
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
func (t *Token) SetId(id Id) error {
	if !t.id.IsValid() {
		return errors.New("token id already set")
	}
	if id.IsValid() {
		return errors.New("cannot set empty token id")
	}
	t.id = id
	return nil
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
	t.updatedAt = DateTime(ti)
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

// touch updates the updatedAt timestamp.
func (u *Token) touch() {
	u.updatedAt = DateTime(time.Now().UTC())
}
