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
}

func (p *SubscriptionPlan) SetId(id int) error {
	newId, err := NewId(id)

	if err != nil {
		return err
	}

	p.id = newId
	return nil
}

func NewSubcriptionPlan(_name, _description string, _duration, _amountInUsd int, maxQuestions int, maxMaterials int, maxUploadSize int, teacherCount int, subType SubscriptionType) (*SubscriptionPlan, error) {
	var errs []error
	id := Id(0) //initially set to zero so persistence layer(repo) can update this via set Id
	name, err := NewName(_name)
	if err != nil {
		errs = append(errs, err)

	}
	description, err := NewDescription(_name)
	if err != nil {
		errs = append(errs, err)

	}
	duration, err := NewDuration(_duration)
	if err != nil {
		errs = append(errs, err)

	}
	price, err := NewPrice(_amountInUsd)
	if err != nil {
		errs = append(errs, err)

	}
	limit, err := NewLimit(maxQuestions, maxMaterials, maxUploadSize, teacherCount)
	if err != nil {
		errs = append(errs, err)

	}
	isTypeValid := subType.IsValid()
	if !isTypeValid {
		errs = append(errs, errors.New("invalid subscription type"))

	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...) // is there a way to retieve individual errors after joining

	}
	return &SubscriptionPlan{id: id, name: name, description: description, duration: duration, price: price, limit: limit, subType: subType}, nil

}

// --- Getter Methods for SubscriptionPlan (making fields private) ---

func (sp *SubscriptionPlan) Id() Id {
	return sp.id
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
