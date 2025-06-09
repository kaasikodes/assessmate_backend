package payment

import "errors"

// TransactionEntityType
type TransactionEntityType string

var (
	SubscriptionEntity TransactionEntityType = "subscription"
)

func NewTransactionEntityType(val string) (TransactionEntityType, error) {
	if isValidTransactionEntityType(val) {
		return TransactionEntityType(val), nil

	}
	return "", errors.New("the transaction entity type is not recognized")

}
func (st TransactionEntityType) String() string {
	return string(st)
}
func (st TransactionEntityType) IsValid() bool {
	return isValidTransactionEntityType(string(st))
}

// IsValid checks if the ProviderType is one of the predefined valid types.
func isValidTransactionEntityType(val string) bool {
	switch TransactionEntityType(val) {
	case SubscriptionEntity:
		return true
	default:
		return false
	}
}
