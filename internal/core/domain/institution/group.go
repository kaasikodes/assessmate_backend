package institution

import (
	"errors"
	"time"
)

type Group struct {
	id                Id
	name              Name
	description       Description
	createdAt         DateTime
	updatedAt         DateTime
	staff             []Staff
	accessibleCourses []Course
}

// NewGroup creates a new Group with the provided name and description.
func NewGroup(name Name, description Description) (*Group, error) {
	if name.IsEmpty() {
		return nil, errors.New("group name cannot be empty")
	}
	if description.IsEmpty() {
		return nil, errors.New("group description cannot be empty")
	}

	now := DateTime(time.Now())

	return &Group{
		name:              name,
		description:       description,
		createdAt:         now,
		updatedAt:         now,
		staff:             []Staff{},
		accessibleCourses: []Course{},
	}, nil
}

// SetId sets the unique identifier of the group.
func (g *Group) SetId(id Id) {
	g.id = id
}

// AddStaff adds a staff member to the group.
func (g *Group) AddStaff(s Staff) {
	g.staff = append(g.staff, s)
	g.touch()
}

// AddCourse grants access to a course for the group.
func (g *Group) AddCourse(c Course) {
	g.accessibleCourses = append(g.accessibleCourses, c)
	g.touch()
}

// UpdateName changes the group's name.
func (g *Group) UpdateName(newName Name) error {
	if newName.IsEmpty() {
		return errors.New("group name cannot be empty")
	}
	g.name = newName
	g.touch()
	return nil
}

// UpdateDescription changes the group's description.
func (g *Group) UpdateDescription(newDesc Description) error {
	if newDesc.IsEmpty() {
		return errors.New("group description cannot be empty")
	}
	g.description = newDesc
	g.touch()
	return nil
}

// touch updates the updatedAt timestamp.
func (g *Group) touch() {
	g.updatedAt = DateTime(time.Now())
}

// ----------- Getters -----------

func (g *Group) Id() Id {
	return g.id
}

func (g *Group) Name() Name {
	return g.name
}

func (g *Group) Description() Description {
	return g.description
}

func (g *Group) CreatedAt() DateTime {
	return g.createdAt
}

func (g *Group) UpdatedAt() DateTime {
	return g.updatedAt
}

func (g *Group) Staff() []Staff {
	return g.staff
}

func (g *Group) AccessibleCourses() []Course {
	return g.accessibleCourses
}
