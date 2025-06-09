package payment

import "errors"

type Transaction struct {
	id                    Id
	providerType          ProviderType
	providerTransactionId ProviderTransactionId
	amount                Amount //in kobo
	paidAt                *DateTime
	transactionEntityType TransactionEntityType
	transactionEntityId   Id
	createdAt             DateTime
	updatedAt             DateTime
	paymentLink           *Url
	meta                  *Meta

	// ProviderTrasactionId
	// Meta
	// Amount
	// CreatedAt
	// UpdatedAt
	// PaidAt
	// TransactionEntityType => subscription
	// TransactionEntityId => 1
}

func (p *Transaction) SetId(id int) error {
	newId, err := NewId(id)

	if err != nil {
		return err
	}

	p.id = newId
	return nil
}

type TransactionParams struct {
	ProviderType          string
	ProviderTransactionId string
	Amount                int //in kobo
	PaidAt                *string
	TransactionEntityType string
	TransactionEntityId   int
	CreatedAt             string
	UpdatedAt             string
}

func NewTransaction(params TransactionParams) (*Transaction, error) {
	var errs []error
	id := Id(0) // Initially set to 0; repo sets actual ID

	providerType, err := NewProviderType(params.ProviderType)
	if err != nil {
		errs = append(errs, err)
	}
	providerTransactionId, err := NewTransactionId(params.ProviderTransactionId)
	if err != nil {
		errs = append(errs, err)
	}
	amount, err := NewAmount(params.Amount)
	if err != nil {
		errs = append(errs, err)
	}

	paidAt, err := NewDateTime(*params.PaidAt)
	if err != nil {
		errs = append(errs, err)
	}
	updatedAt, err := NewDateTime(params.UpdatedAt)
	if err != nil {
		errs = append(errs, err)
	}
	createdAt, err := NewDateTime(params.CreatedAt)
	if err != nil {
		errs = append(errs, err)
	}
	transactionEntityType, err := NewTransactionEntityType(params.TransactionEntityType)
	if err != nil {
		errs = append(errs, err)
	}
	transactionEntityId, err := NewId(params.TransactionEntityId)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Transaction{
		id:                    id,
		providerType:          providerType,
		providerTransactionId: providerTransactionId,
		amount:                amount,
		paidAt:                paidAt,
		transactionEntityType: transactionEntityType,
		transactionEntityId:   transactionEntityId,
		createdAt:             *createdAt,
		updatedAt:             *updatedAt,
	}, nil
}

// --- Getter Methods for SubscriptionPlan (making fields private) ---

func (t *Transaction) Id() Id {
	return t.id
}
func (t *Transaction) ProviderType() ProviderType {
	return t.providerType
}

func (t *Transaction) UpdatedAt() DateTime {
	return t.updatedAt
}
func (t *Transaction) CreatedAt() DateTime {
	return t.createdAt
}
func (t *Transaction) PaidAt() *DateTime {
	return t.paidAt
}
func (t *Transaction) Amount() Amount {
	return t.amount
}
func (t *Transaction) Meta() Meta {
	return *t.meta
}
func (t *Transaction) TransactionEntityId() Id {
	return t.transactionEntityId
}
func (t *Transaction) TransactionEntityType() TransactionEntityType {
	return t.transactionEntityType
}
func (t *Transaction) ProviderTransactionId() ProviderTransactionId {
	return t.providerTransactionId
}
func (t *Transaction) PaymentLink() Url {
	return *t.paymentLink
}
