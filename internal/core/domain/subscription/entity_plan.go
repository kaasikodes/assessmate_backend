package subscription

import "errors"

// TODO: Implement on-chat changes before proceeding, finish core before moving to adapters, so just ports
// TODO: Also learn about bounded context and how and where to use them
// TODO: Better attitude
// TODO: Look into opensource

// entities
type SubscriptionPlan struct {
	id          Id
	name        Name
	description Description
	duration    Duration
	price       Price
	limit       Limit
	subType     SubscriptionType
	isActive    IsActive
}

func (p *SubscriptionPlan) SetId(id int) error {
	newId, err := NewId(id)

	if err != nil {
		return err
	}

	p.id = newId
	return nil
}

type SubscriptionPlanParams struct {
	Name           string
	Description    string
	DurationInDays int
	AmountInUSD    float64
	Limit          LimitParams
	SubType        string
	IsActive       bool
}

type LimitParams struct {
	MaxQuestions  int
	MaxMaterials  int
	MaxUploadSize int
	TeacherCount  int
}

func NewSubscriptionPlan(params SubscriptionPlanParams) (*SubscriptionPlan, error) {
	var errs []error
	id := Id(0) // Initially set to 0; repo sets actual ID

	name, err := NewName(params.Name)
	if err != nil {
		errs = append(errs, err)
	}

	isActive, err := NewIsActive(params.IsActive)
	if err != nil {
		errs = append(errs, err)
	}
	description, err := NewDescription(params.Description)
	if err != nil {
		errs = append(errs, err)
	}

	duration, err := NewDuration(params.DurationInDays)
	if err != nil {
		errs = append(errs, err)
	}

	price, err := NewPrice(params.AmountInUSD)
	if err != nil {
		errs = append(errs, err)
	}

	limit, err := NewLimit(
		params.Limit.MaxQuestions,
		params.Limit.MaxMaterials,
		params.Limit.MaxUploadSize,
		params.Limit.TeacherCount,
	)
	if err != nil {
		errs = append(errs, err)
	}
	subType, err := NewSubcriptionType(params.SubType)

	if !subType.IsValid() {
		errs = append(errs, err, errors.New("invalid subscription type"))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &SubscriptionPlan{
		id:          id,
		name:        name,
		description: description,
		duration:    duration,
		price:       price,
		limit:       limit,
		subType:     subType,
		isActive:    isActive,
	}, nil
}

// --- Getter Methods for SubscriptionPlan (making fields private) ---

func (sp *SubscriptionPlan) Id() Id {
	return sp.id
}
func (sp *SubscriptionPlan) IsActive() IsActive {
	return sp.isActive
}

func (sp *SubscriptionPlan) Name() Name {
	return sp.name
}

func (sp *SubscriptionPlan) Description() Description {
	return sp.description
}

func (sp *SubscriptionPlan) Duration() Duration {
	return sp.duration
}

func (sp *SubscriptionPlan) Price() Price {
	return sp.price
}

func (sp *SubscriptionPlan) Limit() Limit {
	return sp.limit
}

func (sp *SubscriptionPlan) SubscriptionType() SubscriptionType {
	return sp.subType
}
