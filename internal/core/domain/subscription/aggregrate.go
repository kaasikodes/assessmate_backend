package subscription

import (
	"errors"
	"fmt"
	"strings"
)

// aggregrate
type SubscriptionPlanAggregrate struct {
	Plan         SubscriptionPlan
	Institutions []Institution
}

type Subscription struct {
	Id              Id
	Plan            SubscriptionPlan
	UserId          Id //all subscriptions are tied to user, if insittution or enterprise, there will be subscription_institution to make relatioship, and a IsActive is the subscription that is active
	Subscriber      Subscriber
	HasPaid         HasPaid
	ExpiresAt       DateTime
	PaidAt          *DateTime
	TransactionId   *TransactionId
	PaymentProvider *PaymentProvider
	Meta            *Meta
}
type PaymentProvider string

var (
	Paystack PaymentProvider = "paystack"
	Mono     PaymentProvider = "mono"
	KorahPay PaymentProvider = "korahpay"
)

func NewPaymentProvider(val string) (PaymentProvider, error) {
	val = strings.ToLower(strings.TrimSpace(val))

	switch PaymentProvider(val) {
	case Paystack, Mono, KorahPay:
		return PaymentProvider(val), nil
	default:
		return "", errors.New("invalid payment provider: " + val)
	}
}

type TransactionId string

func NewTransactionId(val string) (TransactionId, error) {
	return TransactionId(val), nil
}

type Meta map[string]string

func NewMeta() (*Meta, error) {
	meta := make(Meta)
	return &meta, nil

}

// Add inserts a key-value pair if the key does not already exist.
func (m *Meta) Add(key, value string) error {
	if m == nil {
		return errors.New("meta is nil")
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("key cannot be empty")
	}

	if _, exists := (*m)[key]; exists {
		return fmt.Errorf("key %q already exists in meta", key)
	}

	(*m)[key] = value
	return nil
}

// Get returns the value associated with a key.
func (m *Meta) Get(key string) (string, bool) {
	if m == nil {
		return "", false
	}

	value, exists := (*m)[key]
	return value, exists
}

// Remove deletes a key from the map if it exists.
func (m *Meta) Remove(key string) error {
	if m == nil {
		return errors.New("meta is nil")
	}

	key = strings.TrimSpace(key)
	if _, exists := (*m)[key]; !exists {
		return fmt.Errorf("key %q does not exist in meta", key)
	}

	delete(*m, key)
	return nil
}

type Subscriber struct {
	Id            Id
	Name          Name
	Email         Email
	Subscriptions []Subscription
}

// external entities
type Institution struct {
	Id    int
	Email string
	Name  string
	Owner User
}

type User struct {
	Id int
}
