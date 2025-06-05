package subscription

import (
	"context"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/subscription"
)

// ports should conform to language of core(in this case the domain and application)
type SubRepository interface {
	CreatePlan(ctx context.Context, payload *subscription.SubscriptionPlan) (*subscription.SubscriptionPlan, error)
	UpdatePlan(ctx context.Context, payload *subscription.SubscriptionPlan) (*subscription.SubscriptionPlan, error)
	FindPlanById(ctx context.Context, id int) (*subscription.SubscriptionPlan, error)
	DeletePlan(ctx context.Context, id int) error
}
