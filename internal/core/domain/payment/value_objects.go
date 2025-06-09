package payment

import (
	"errors"
	"net/url"
	"time"
)

// Transactions
// Url
type Url string

// NewUrl validates and returns a Url type
func NewUrl(val string) (Url, error) {
	parsed, err := url.ParseRequestURI(val)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("invalid URL")
	}
	return Url(val), nil
}

// String returns the string representation of the Url
func (u Url) String() string {
	return string(u)
}

// Id
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

// ProviderType
type ProviderType string

var (
	Paystack ProviderType = "paystack"
	Mono     ProviderType = "mono"
	KorahPay ProviderType = "korahpay"
)

func NewProviderType(val string) (ProviderType, error) {
	if isValidProviderType(val) {
		return ProviderType(val), nil

	}
	return "", errors.New("the provider type is not recognized")

}
func (st ProviderType) IsValid() bool {
	return isValidProviderType(string(st))
}
func (st ProviderType) String() string {
	return string(st)
}

// IsValid checks if the ProviderType is one of the predefined valid types.
func isValidProviderType(val string) bool {
	switch ProviderType(val) {
	case Paystack, Mono, KorahPay:
		return true
	default:
		return false
	}
}

type ProviderTransactionId string

func NewTransactionId(val string) (ProviderTransactionId, error) {
	if val == "" {
		return "", errors.New("transaction Id cannot be empty")
	}
	return ProviderTransactionId(val), nil

}

type Amount int // in kobo
func NewAmount(value int) (Amount, error) {
	if value <= 10_000 {
		return 0, errors.New("amount has to be greater or equal to 1000")
	}
	return Amount(value), nil
}

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

// Amount
// CreatedAt
// UpdatedAt
// Meta
// PaidAt
// ProviderTransactionId
// TransactionEntityType => subscription
// TransactionEntityId => 1
