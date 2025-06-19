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
	password   password
	createdAt  DateTime
	updatedAt  DateTime
	verifiedAt *DateTime
	deletedAt  *DateTime
}

// NewUser creates a new user instance with default status and timestamps.
func NewUser(name Name, email Email) (*User, error) {

	if name.IsEmpty() {
		return nil, errors.New("user name cannot be empty")
	}
	if email.IsEmpty() {
		return nil, errors.New("user email cannot be empty")
	}

	now := DateTime(time.Now().UTC())

	return &User{
		name:      name,
		status:    Active,
		email:     email,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// SetId sets the user's ID if not already set.
func (i *User) SetId(id Id) {
	i.id = id
}

// SetPassword
func (u *User) SetPasswordHash(hash []byte) error {
	u.password.SetHash(hash)
	u.touch()
	return nil
}

// SetPassword
func (u *User) SetPassword(text string) error {
	err := u.password.NewHash(text)
	if err != nil {
		return err
	}
	u.touch()
	return nil
}

// ComparePassword
func (u *User) ComparePassword(text string) bool {
	return u.password.Compare(text)

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
func (u *User) GetEmail() Email {
	return u.email
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
func (u *User) PasswordHash() []byte {
	return u.password.hash
}

// touch updates the updatedAt timestamp.
func (u *User) touch() {
	u.updatedAt = DateTime(time.Now().UTC())
}
