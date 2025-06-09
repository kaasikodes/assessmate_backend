package institution

import (
	"errors"
	"time"
)

// Institution represents an institution entity in the domain.
type Institution struct {
	id          Id
	name        Name
	description Description
	email       Email
	staff       []Staff
	createdAt   DateTime
	updatedAt   DateTime
}

// NewInstitution is a factory function that creates a new Institution.
func NewInstitution(name Name, description Description, email Email) (*Institution, error) {
	if name.IsEmpty() {
		return nil, errors.New("name cannot be empty")
	}
	if email.IsEmpty() {
		return nil, errors.New("email cannot be empty")
	}

	now := DateTime(time.Now().UTC())
	return &Institution{
		name:        name,
		description: description,
		email:       email,
		staff:       []Staff{},
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// SetId sets the ID of the institution. Usually used when loaded from persistence.
func (i *Institution) SetId(id Id) {
	i.id = id
}

// AddStaff adds a staff member to the institution.
func (i *Institution) AddStaff(s Staff) {
	i.staff = append(i.staff, s)
	i.touch()
}

// Getters
func (i *Institution) Id() Id {
	return i.id
}

func (i *Institution) Name() Name {
	return i.name
}

func (i *Institution) Description() Description {
	return i.description
}

func (i *Institution) Email() Email {
	return i.email
}

func (i *Institution) Staff() []Staff {
	return i.staff
}

func (i *Institution) CreatedAt() DateTime {
	return i.createdAt
}

func (i *Institution) UpdatedAt() DateTime {
	return i.updatedAt
}

// UpdateName updates the institution name and updatedAt timestamp.
func (i *Institution) UpdateName(newName Name) error {
	if newName.IsEmpty() {
		return errors.New("name cannot be empty")
	}
	i.name = newName
	i.touch()
	return nil
}

// touch updates the updatedAt timestamp.
func (i *Institution) touch() {
	i.updatedAt = DateTime(time.Now().UTC())
}
