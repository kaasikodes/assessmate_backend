package institution

import (
	"errors"
	"time"
)

type Course struct {
	id          Id
	name        Name
	description Description
	createdAt   DateTime
	updatedAt   DateTime
}

// NewCourse creates a new course instance with validation.
// It sets createdAt and updatedAt to the current time.
func NewCourse(name Name, description Description) (*Course, error) {
	if name.IsEmpty() {
		return nil, errors.New("course name cannot be empty")
	}
	if description.IsEmpty() {
		return nil, errors.New("course description cannot be empty")
	}

	now := DateTime(time.Now())

	return &Course{
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// SetId assigns an ID to the course (useful post-persistence).
func (c *Course) SetId(id Id) {
	c.id = id
}

// UpdateName changes the course name and updates the updatedAt timestamp.
func (c *Course) UpdateName(newName Name) error {
	if newName.IsEmpty() {
		return errors.New("course name cannot be empty")
	}
	c.name = newName
	c.touch()
	return nil
}

// UpdateDescription changes the course description and updates the updatedAt timestamp.
func (c *Course) UpdateDescription(newDescription Description) error {
	if newDescription.IsEmpty() {
		return errors.New("course description cannot be empty")
	}
	c.description = newDescription
	c.touch()
	return nil
}

// Internal helper to update updatedAt field.
func (c *Course) touch() {
	c.updatedAt = DateTime(time.Now())
}

// ----------- Getters -----------

func (c *Course) Id() Id {
	return c.id
}

func (c *Course) Name() Name {
	return c.name
}

func (c *Course) Description() Description {
	return c.description
}

func (c *Course) CreatedAt() DateTime {
	return c.createdAt
}

func (c *Course) UpdatedAt() DateTime {
	return c.updatedAt
}
