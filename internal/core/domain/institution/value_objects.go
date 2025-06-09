package institution

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// id
type Id int

func NewId(id int) (Id, error) {
	if id == 0 {
		return 0, errors.New("id cannot be zero")

	}
	if id < 0 {
		return 0, errors.New("id cannot be negative")
	}

	return Id(id), nil
}
func (i Id) Value() int {
	return int(i)
}

// name
type Name string

func (n Name) String() string {
	return string(n)

}

func (n Name) IsEmpty() bool {
	return strings.TrimSpace(string(n)) == ""

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

// description
type Description string

func (n Description) String() string {
	return string(n)

}
func (n Description) IsEmpty() bool {
	return strings.TrimSpace(string(n)) == ""

}
func NewDescription(description string) (Description, error) {
	// cannot be empty
	description = strings.TrimSpace(description)
	if description == "" {
		return "", errors.New("description cannot be empty")
	}
	minLength := 150
	maxLength := 800
	// must be at least min characters long
	if len(description) < minLength {
		return "", (fmt.Errorf("description must be at least %d charaters long", minLength))

	}
	// must be at least max characters long
	if len(description) > maxLength {
		return "", (fmt.Errorf("description must not exceed %d charaters", maxLength))

	}

	return Description(description), nil

}

// email
type Email string

func (n Email) String() string {
	return string(n)

}
func (n Email) IsEmpty() bool {
	return strings.TrimSpace(string(n)) == ""

}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return "", errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return "", errors.New("invalid email format")
	}

	return Email(email), nil
}

// datetime
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

// staff status
type StaffStatus string

var (
	Active   StaffStatus = "active"
	InActive StaffStatus = "inactive"
)

func NewStaffStatus(val string) (StaffStatus, error) {
	if isValidStaffStatus(val) {
		return StaffStatus(val), nil

	}
	return "", errors.New("the staff status is not recognized")

}
func (st StaffStatus) IsValid() bool {
	return isValidStaffStatus(string(st))
}
func (st StaffStatus) String() string {
	return string(st)
}

// IsValid checks if the StaffStatus is one of the predefined valid types.
func isValidStaffStatus(val string) bool {
	switch StaffStatus(val) {
	case Active, InActive:
		return true
	default:
		return false
	}
}
