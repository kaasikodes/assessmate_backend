package subscription

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// value objects
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

// This is a convenient getter method.
func (i Id) Value() int {
	return int(i)
}

// String implements the fmt.Stringer interface for easy printing.
func (i Id) String() string {
	return strconv.Itoa(int(i))
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

type Description string

func (n Description) String() string {
	return string(n)

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

type Duration time.Duration

func NewDuration(durationInSecs int) (Duration, error) {

	// duration must be at least 2 weeks
	// duration should not exceed 4 months
	const (
		// Define minimum and maximum durations using time.Duration constants for clarity
		minDuration = 2 * 7 * 24 * time.Hour  // 2 weeks in hours
		maxDuration = 4 * 30 * 24 * time.Hour // Approximately 4 months (using 30 days/month for simplicity)
		// Note: For precise month calculations, you might need to use time.AddDate
		// or consider the actual days in a month if your domain requires that precision.
		// For a typical "duration" value object, a fixed number of days per month is often acceptable.
	)

	// Convert the input int seconds to a time.Duration
	d := time.Duration(durationInSecs) * time.Second

	// Duration must be at least 2 weeks
	if d < minDuration {
		return 0, fmt.Errorf("duration must be at least 2 weeks (%.0f seconds)", minDuration.Seconds())
	}

	// Duration should not exceed 4 months
	if d > maxDuration {
		return 0, fmt.Errorf("duration must not exceed 4 months (%.0f seconds)", maxDuration.Seconds())
	}

	// Return the Duration directly, not a pointer.
	// Since time.Duration is a value type, our Duration alias is also a value type.
	return Duration(d), nil

}

func (d Duration) Seconds() float64 {
	return time.Duration(d).Seconds()
}

func (d Duration) String() string {
	return time.Duration(d).String()
}

type Price struct {
	amount   float64
	currency Currency
}

func NewPrice(amount int) (Price, error) {
	// check amount should not be negative
	// amount should not exceed 10_000_000
	// all price when set have a currency of USD
	const (
		minAmount = 0          // Amount should not be negative
		maxAmount = 10_000_000 // 10,000,000 (assuming base unit like cents for int, or dollars for float)
	)

	// Ensure amount is not negative
	if amount < minAmount {
		return Price{}, errors.New("price amount cannot be negative")
	}

	// Ensure amount does not exceed the maximum
	if amount > maxAmount {
		return Price{}, fmt.Errorf("price amount cannot exceed %d", maxAmount)
	}

	// All prices when set have a currency of USD by default
	// We'll store the float64 representation directly, assuming the `amount int`
	// was passed in the smallest unit (e.g., cents) or adjust if it's already dollars.
	// For simplicity, let's assume `amount` is in the primary unit (e.g., dollars)
	// and we're just directly using it as float64.
	// If `amount` is in cents, you'd do: `float64(amount) / 100.0`
	return Price{
		amount:   float64(amount), // Assuming amount is in primary units (e.g., dollars)
		currency: USD,
	}, nil
}

// Convert converts the price to a specified target currency.
// It returns a new Price object with the converted amount and target currency.
// It returns an error if conversion is not possible or rates are missing.
func (p Price) Convert(targetCurrency Currency) (Price, error) {
	if p.currency == targetCurrency {
		return p, nil // No conversion needed
	}

	// Use the currency package's conversion function
	convertedAmount, err := ConvertCurrency(p.amount, p.currency, targetCurrency)
	if err != nil {
		return Price{}, fmt.Errorf("failed to convert price from %s to %s: %w", p.currency, targetCurrency, err)
	}

	// Round to 2 decimal places for typical currency precision, or more if needed
	// Use math.Round(x * 100) / 100 for proper rounding
	roundedAmount := math.Round(convertedAmount*100) / 100

	return Price{
		amount:   roundedAmount,
		currency: targetCurrency,
	}, nil
}

// Amount returns the numerical value of the price.
func (p Price) Amount() float64 {
	return p.amount
}

// Currency returns the currency unit of the price.
func (p Price) Currency() Currency {
	return p.currency
}

// SubscriptionType represents a predefined type of subscription.
// It is a string-based Value Object, identified by its literal value.
type SubscriptionType string

// Predefined constants for valid SubscriptionType values.
const (
	SubscriptionTypePersonal    SubscriptionType = "personal"
	SubscriptionTypeInstitution SubscriptionType = "institution"
)

// IsValid checks if the SubscriptionType is one of the predefined valid types.
func (st SubscriptionType) IsValid() bool {
	switch st {
	case SubscriptionTypePersonal, SubscriptionTypeInstitution:
		return true
	default:
		return false
	}
}

// Limit defines the usage quotas for a given subscription type.
// It is an immutable Value Object, encapsulating various numerical limits.
type Limit struct {
	MaxQuestions  int // Maximum number of questions allowed
	MaxMaterials  int // Maximum number of learning materials/documents
	MaxUploadSize int // Maximum total upload size in megabytes (MB)
	TeacherCount  int // Number of teachers allowed (relevant for institution/enterprise)
}

// NewLimit creates a new Limit value object.
// This constructor allows you to define specific limits.
// All limits must be non-negative. It returns a Limit and an error if validation fails.
func NewLimit(maxQuestions, maxMaterials, maxUploadSize, teacherCount int) (Limit, error) {
	if maxQuestions < 0 {
		return Limit{}, errors.New("MaxQuestions cannot be negative")
	}
	if maxMaterials < 0 {
		return Limit{}, errors.New("MaxMaterials cannot be negative")
	}
	if maxUploadSize < 0 {
		return Limit{}, errors.New("MaxUploadSize cannot be negative")
	}
	if teacherCount < 0 {
		return Limit{}, errors.New("TeacherCount cannot be negative")
	}
	// not to be exceeded values
	if maxQuestions > 40 {
		return Limit{}, fmt.Errorf("MaxQuestions cannot exceed %d", 40)
	}
	if maxMaterials > 3 {
		return Limit{}, fmt.Errorf("MaxMaterials cannot exceed %d", 3)
	}
	if maxUploadSize > 5 {
		return Limit{}, fmt.Errorf("MaxUploadSize cannot exceed %d", 5)
	}
	if teacherCount > 200 {
		return Limit{}, fmt.Errorf("TeacherCount cannot exceed %d", 200)
	}

	return Limit{
		MaxQuestions:  maxQuestions,
		MaxMaterials:  maxMaterials,
		MaxUploadSize: maxUploadSize,
		TeacherCount:  teacherCount,
	}, nil
}
