package payment

import (
	"context"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/payment"
)

// ports should conform to language of core(in this case the domain and not application, as application is a bridge for adapter to domain(business) logic)
type PaymentRepository interface {
	InitiatePayment(ctx context.Context, payload *payment.Transaction) (*payment.Transaction, error)
	MarkTransactionAsPaid(ctx context.Context, transactionId payment.Id) (*payment.Transaction, error)
}
