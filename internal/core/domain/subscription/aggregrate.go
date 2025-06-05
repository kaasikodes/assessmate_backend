package subscription

// aggregrate
type SubscriptionPlanAggregrate struct {
	Plan         SubscriptionPlan
	Institutions []Institution
}

type Subscription struct {
	Id           int
	PlanId       int
	SubscriberId int //depends on the plan type if inistution then will be an insttution id, personal will be a userId
	Subscriber   Subscriber
}
type Subscriber struct {
	Id   int
	Name string
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
