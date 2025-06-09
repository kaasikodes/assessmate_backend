package subscription

import (
	"context"
	"errors"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/subscription"
)

var ErrPlanNotFound = errors.New("subscription plan not found")

// ports should conform to language of core(in this case the domain and not application, as application is a bridge for adapter to domain(business) logic)
type SubscriptionRepository interface {
	CreatePlan(ctx context.Context, payload *subscription.SubscriptionPlan) (*subscription.SubscriptionPlan, error)
	DeletePlan(ctx context.Context, id subscription.Id) error
	FindPlanById(ctx context.Context, id subscription.Id) (*subscription.SubscriptionPlan, error)
	GetSubscribers(cts context.Context, filter *subscription.SubscriberFilterParams) (result []subscription.Subscriber, total int, err error)
	GetPlans(ctx context.Context, filter *subscription.PlanFilterParams) (result []subscription.SubscriptionPlan, total int, err error)
	UpdatePlan(ctx context.Context, payload *subscription.SubscriptionPlan) (*subscription.SubscriptionPlan, error)
	ActivateOrDeactivatePlan(ctx context.Context, id subscription.Id, isActive subscription.IsActive) (*subscription.SubscriptionPlan, error)
	CreateSubscription(ctx context.Context, planId subscription.Id, userId subscription.Id) (*subscription.Subscriber, error)
}
