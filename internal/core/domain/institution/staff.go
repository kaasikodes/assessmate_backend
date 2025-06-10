package institution

import (
	"errors"
	"time"
)

// Staff represents a staff member within an institution.
type Staff struct {
	id        Id
	name      Name
	email     Email
	status    StaffStatus
	deletedAt *DateTime
	createdAt DateTime
	updatedAt DateTime
}

// NewStaff creates a new Staff domain entity.
// Note: status defaults to StaffActive if empty.
func NewStaff(name Name, email Email, status StaffStatus) (*Staff, error) {
	if name.IsEmpty() {
		return nil, errors.New("staff name cannot be empty")
	}
	if email.IsEmpty() {
		return nil, errors.New("staff email cannot be empty")
	}
	if status == "" {
		status = InActive
	}
	return &Staff{
		name:   name,
		email:  email,
		status: status,
	}, nil
}

// SetId assigns a persisted Id once retrieved from storage.
func (s *Staff) SetId(id Id) {
	s.id = id
}

// MarkDeleted records soft‑deletion with a timestamp.
func (s *Staff) MarkDeleted(t time.Time) {
	dt := DateTime(t)
	s.deletedAt = &dt
}

// ------------- Getters -------------

func (s *Staff) Id() Id {
	return s.id
}
func (s *Staff) Name() Name {
	return s.name
}
func (s *Staff) Email() Email {
	return s.email
}
func (s *Staff) Status() StaffStatus {
	return s.status
}
func (s *Staff) DeletedAt() *DateTime {
	return s.deletedAt
}
func (s *Staff) CreatedAt() DateTime {
	return s.createdAt
}
func (s *Staff) UpdatedAt() DateTime {
	return s.updatedAt
}
func (s *Staff) IsDeleted() bool {
	return s.deletedAt != nil
}
