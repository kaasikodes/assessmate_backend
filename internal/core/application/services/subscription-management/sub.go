package subscriptionmanagement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/subscription"
	sub_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/subscription"
)

type (
	LimitPayload struct {
		MaxQuestions  int
		MaxMaterials  int
		MaxUploadSize int
		TeacherCount  int
	}

	CreatePlanPayload struct {
		Name           string
		Description    string
		DurationInDays int
		PriceInUsd     float64
		Limit          LimitPayload
		Type           string
	}
	CreatePlanResponse struct {
		Id             int
		Name           string
		Description    string
		DurationInDays int
		PriceInUsd     float64
		Limit          LimitPayload
		Type           string
	}
	SubscriberCount struct {
		Active  int // expiry date is a day after present day or in the future
		Expired int // expiry date is in the past
		Due     int // expiry date is the present day
	}
	GetPlanResponse struct {
		Id              int
		Name            string
		Description     string
		DurationInDays  int
		PriceInUsd      float64
		Limit           LimitPayload
		Type            string
		SubscriberCount SubscriberCount
		IsActive        bool
	}
	GetSubscribersResponse struct {
		Total  int
		Result []Subscriber
	}
	Subscriber struct {
		Id                 int
		Name               string
		Email              string
		ActiveSubscription Subscription
		Subscriptions      []Subscription
	}
	Subscription struct {
		Id              int
		HasPaid         bool
		ExpiresAt       time.Time
		PaidAt          *time.Time
		TransactionId   *string
		PaymentProvider *string
		Meta            *map[string]string
	}

	SubscriberFilterParams struct {
		PlanId         *int //TODO: Make Private to adhere to ddd, create Getter Methods
		ExpiryDateFrom string
		ExpiryDateTo   string
		HasPaid        *bool

		// SubscriptionType *SubscriptionType //TODO: Add ability to  filter by suvscription type
	}
)

type SubscriptionManagementService struct {
	planRepo sub_repo.SubscriptionRepository
}

// create plan
func (s *SubscriptionManagementService) CreatePlan(ctx context.Context, payload CreatePlanPayload) (*CreatePlanResponse, error) {
	// converts payload to domain  language - verifies input and coverts to domain language
	plan, err := subscription.NewSubscriptionPlan(subscription.SubscriptionPlanParams{
		Name:           payload.Name,
		Description:    payload.Description,
		DurationInDays: payload.DurationInDays,
		AmountInUSD:    (payload.PriceInUsd),
		Limit: subscription.LimitParams{
			MaxQuestions:  payload.Limit.MaxQuestions,
			MaxMaterials:  payload.Limit.MaxMaterials,
			MaxUploadSize: payload.Limit.MaxUploadSize,
			TeacherCount:  payload.Limit.TeacherCount,
		},
		IsActive: false,
	})
	//if there is an issue with data let us know
	if err != nil {
		return nil, err
	}
	// the repo persists the data because it trusts that the data has been validated, and its job is to persist data and not validate data
	plan, err = s.planRepo.CreatePlan(ctx, plan)
	//if there is an issue persisting data let us know (is db down, is there a concurrency issue the data version is above the request)
	if err != nil {
		return nil, err
	}

	return &CreatePlanResponse{
		Id:             plan.Id().Value(),
		Name:           plan.Name().String(),
		Description:    plan.Description().String(),
		DurationInDays: int(plan.Duration().Days()),
		PriceInUsd:     (plan.Price().Amount()),
		Limit: LimitPayload{
			MaxQuestions:  plan.Limit().MaxQuestions,
			MaxMaterials:  plan.Limit().MaxMaterials,
			MaxUploadSize: plan.Limit().MaxUploadSize,
			TeacherCount:  plan.Limit().TeacherCount,
		},
		Type: string(plan.SubscriptionType()),
	}, nil

}

// get plan by id
func (s *SubscriptionManagementService) GetPlan(ctx context.Context, planId int) (*GetPlanResponse, error) {
	id, err := subscription.NewId(planId)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id: %w", err)
	}
	// Rerieve plan from persistent storage
	plan, err := s.planRepo.FindPlanById(ctx, id)
	if err != nil {
		if errors.Is(err, sub_repo.ErrPlanNotFound) {
			return nil, fmt.Errorf("cannot retrieve plan: %w", err)
		}
		return nil, fmt.Errorf("error checking if plan exists: %w", err)
	}
	hasPaid := true
	// 	Active   int // expiry date is a day after present day or in the future
	active, err := s.GetSubscribersList(ctx, &SubscriberFilterParams{PlanId: &planId, ExpiryDateFrom: time.Now().Add(24 * time.Hour).Format(time.RFC3339), ExpiryDateTo: "", HasPaid: &hasPaid})
	if err != nil {
		return nil, err
	}
	activeSubscribers := active.Total

	// Expired: expiry date is in the past
	expired, err := s.GetSubscribersList(ctx, &SubscriberFilterParams{PlanId: &planId, ExpiryDateTo: time.Now().Format(time.RFC3339), HasPaid: &hasPaid})
	if err != nil {
		return nil, err
	}
	expiredSubscribers := expired.Total

	// Due: Subscriptions with expiry date between 5 days before now and end of today
	startOfDueDay := time.Now().Truncate(24 * time.Hour).Add(-5 * 24 * time.Hour)         // 5 days before midnight today
	endOfToday := time.Now().Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond) // End of today (23:59:59.999...)

	due, err := s.GetSubscribersList(ctx, &SubscriberFilterParams{
		PlanId:         &planId,
		ExpiryDateFrom: startOfDueDay.Format(time.RFC3339),
		ExpiryDateTo:   endOfToday.Format(time.RFC3339),
		HasPaid:        &hasPaid,
	})
	if err != nil {
		return nil, err
	}
	dueSubscribers := due.Total

	return &GetPlanResponse{
		Id:             plan.Id().Value(),
		Name:           plan.Id().String(),
		Description:    plan.Description().String(),
		DurationInDays: int(plan.Duration().Days()),
		PriceInUsd:     plan.Price().Amount(),
		Limit: LimitPayload{
			MaxQuestions:  plan.Limit().MaxQuestions,
			MaxMaterials:  plan.Limit().MaxMaterials,
			MaxUploadSize: plan.Limit().MaxUploadSize,
			TeacherCount:  plan.Limit().TeacherCount,
		},
		Type:     string(plan.SubscriptionType()),
		IsActive: plan.IsActive().Bool(),
		SubscriberCount: SubscriberCount{
			Active:  activeSubscribers,
			Expired: expiredSubscribers,
			Due:     dueSubscribers,
		}, //TODO: Get Subscribers and the count and append here once the subscriber service is created
	}, nil

}

//get list of plans (filter by is Active)

// get list of subscribers (and be able to filter by plan, exiryDate, hasPaid)
func (s *SubscriptionManagementService) GetSubscribersList(ctx context.Context, _filter *SubscriberFilterParams) (*GetSubscribersResponse, error) {
	var filter subscription.SubscriberFilterParams
	if _filter != nil {
		filterParsed, err := subscription.NewSubscriberFilterParams(_filter.PlanId, _filter.ExpiryDateFrom, _filter.ExpiryDateTo, _filter.HasPaid)
		if err != nil {
			return nil, fmt.Errorf("invalid filter params: %w", err)
		}
		filter = *filterParsed
	}

	result, total, err := s.planRepo.GetSubscribers(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subscribers: %w", err)
	}

	subscribers := make([]Subscriber, len(result))
	for i, sub := range result {
		subscribers[i].Id = sub.Id.Value()
		subscribers[i].Name = sub.Name.String()
		subscribers[i].Email = sub.Email.String()
		subscribers[i].Subscriptions = make([]Subscription, len(sub.Subscriptions))

		var activeSub *Subscription
		var earliestExpiry *time.Time

		for j, resSub := range sub.Subscriptions {
			subDTO := Subscription{
				Id:              resSub.Id.Value(),
				ExpiresAt:       resSub.ExpiresAt,
				HasPaid:         resSub.HasPaid.Bool(),
				PaidAt:          resSub.PaidAt,
				PaymentProvider: nil,
				TransactionId:   nil,
				Meta:            nil,
			}

			if resSub.PaymentProvider != nil {
				pp := string(*resSub.PaymentProvider)
				subDTO.PaymentProvider = &pp
			}

			if resSub.TransactionId != nil {
				tid := string(*resSub.TransactionId)
				subDTO.TransactionId = &tid
			}

			if resSub.Meta != nil {
				m := map[string]string(*resSub.Meta)
				subDTO.Meta = &m
			}

			subscribers[i].Subscriptions[j] = subDTO

			// Determine the active subscription: which is the paid subscription with the earliest future expiry date
			if subDTO.HasPaid && subDTO.ExpiresAt.After(time.Now()) {
				if earliestExpiry == nil || subDTO.ExpiresAt.Before(*earliestExpiry) {
					earliestExpiry = &subDTO.ExpiresAt
					activeSub = &subDTO
				}
			}
		}

		subscribers[i].ActiveSubscription = *activeSub
	}

	return &GetSubscribersResponse{
		Total:  total,
		Result: subscribers,
	}, nil
}

func (s *SubscriptionManagementService) verifyPlanExists(ctx context.Context, planId int) (*GetPlanResponse, *subscription.Id, error) {
	id, err := subscription.NewId(planId)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid plan id: %w", err)
	}
	// Ensure plan exists before trying to delete
	plan, err := s.GetPlan(ctx, int(id))
	if err != nil {
		return nil, nil, fmt.Errorf("error checking if plan exists: %w", err)
	}
	return plan, &id, nil
}

// delete plan
func (s *SubscriptionManagementService) DeletePlan(ctx context.Context, planId int) error {
	_, parsedPlanId, err := s.verifyPlanExists(ctx, planId)
	if err != nil {
		return err
	}
	//Delete from persistent storage
	err = s.planRepo.DeletePlan(ctx, *parsedPlanId)
	if err != nil {
		return fmt.Errorf("failed to delete plan with id %d: %w", planId, err)
	}
	return nil

}

// activate/deactivate plan
// subscribeUserToPlan
// update plan
//list of plans
// get single plan
// cancelSubscription
// renewSubscription
// payFor Subscription
//autoRenewCurrentSubscription
