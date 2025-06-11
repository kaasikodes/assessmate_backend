package user

import (
	"errors"
	"time"
)

type User struct {
	id         Id
	name       Name
	email      Email
	status     UserStatus
	createdAt  DateTime
	updatedAt  DateTime
	verifiedAt *DateTime
	deletedAt  *DateTime
}

// NewUser creates a new user instance with default status and timestamps.
func NewUser(id Id, name Name, email Email) (*User, error) {

	if name.IsEmpty() {
		return nil, errors.New("user name cannot be empty")
	}
	if email.IsEmpty() {
		return nil, errors.New("user email cannot be empty")
	}

	now := DateTime(time.Now().UTC())

	return &User{
		id:        id,
		name:      name,
		status:    Active,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// SetId sets the user's ID if not already set.
func (i *User) SetId(id Id) {
	i.id = id
}

// SetVerifiedAt marks the user as verified.
func (u *User) SetVerifiedAt(t time.Time) {
	dt := DateTime(t)
	u.verifiedAt = &dt
	u.touch()
}

// SetDeletedAt marks the user as deleted.
func (u *User) SetDeletedAt(t time.Time) {
	dt := DateTime(t)
	u.deletedAt = &dt
	u.touch()
}

// SetStatus changes the user’s status and updates the timestamp.
func (u *User) SetStatus(status UserStatus) {
	u.status = status
	u.touch()
}

// SetEmail changes the user’s email and updates the timestamp.
func (u *User) SetEmail(email Email) {
	u.email = email
	u.touch()
}

// Getters
func (u *User) GetId() Id {
	return u.id
}

func (u *User) GetName() Name {
	return u.name
}

func (u *User) GetStatus() UserStatus {
	return u.status
}

func (u *User) GetCreatedAt() DateTime {
	return u.createdAt
}

func (u *User) GetUpdatedAt() DateTime {
	return u.updatedAt
}

func (u *User) GetVerifiedAt() *DateTime {
	return u.verifiedAt
}

func (u *User) GetDeletedAt() *DateTime {
	return u.deletedAt
}
func (u *User) IsVerified() bool {
	return u.verifiedAt != nil
}

// touch updates the updatedAt timestamp.
func (u *User) touch() {
	u.updatedAt = DateTime(time.Now().UTC())
}
