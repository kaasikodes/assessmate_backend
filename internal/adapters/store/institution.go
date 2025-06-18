package store

import (
	"context"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/subscription"
)

func (r *MySqlRepo) CreatePlan(ctx context.Context, payload *subscription.SubscriptionPlan) (*subscription.SubscriptionPlan, error)
func (r *MySqlRepo) DeletePlan(ctx context.Context, id subscription.Id) error
func (r *MySqlRepo) FindPlanById(ctx context.Context, id subscription.Id) (*subscription.SubscriptionPlan, error)
func (r *MySqlRepo) GetSubscribers(cts context.Context, filter *subscription.SubscriberFilterParams) (result []subscription.Subscriber, total int, err error)
func (r *MySqlRepo) GetPlans(ctx context.Context, filter *subscription.PlanFilterParams) (result []subscription.SubscriptionPlan, total int, err error)
func (r *MySqlRepo) UpdatePlan(ctx context.Context, payload *subscription.SubscriptionPlan) (*subscription.SubscriptionPlan, error)
func (r *MySqlRepo) ActivateOrDeactivatePlan(ctx context.Context, id subscription.Id, isActive subscription.IsActive) (*subscription.SubscriptionPlan, error)
func (r *MySqlRepo) CreateSubscription(ctx context.Context, planId subscription.Id, userId subscription.Id) (*subscription.Subscriber, error)
