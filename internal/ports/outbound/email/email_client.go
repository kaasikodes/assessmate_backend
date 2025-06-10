package email

import "context"

type MailInput struct {
	Name     string
	Email    string
	Message  string
	Template string
}
type EmailClient interface {
	SendMultiple(ctx context.Context, payload []MailInput) error
	Send(ctx context.Context, payload MailInput) error
}
