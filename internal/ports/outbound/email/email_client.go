package email

import (
	"context"
	"time"
)

type MailInput struct {
	Name     string
	Email    string
	Message  string
	Template string
}

type Notification struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`  // can be null
	IsRead    bool       `json:"isRead"` //defaults to false
	ReadAt    *time.Time `json:"readAt"` //can be null
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
type EmailClient interface {
	Send(ctx context.Context, notification *Notification) error
	SendMultiple(ctx context.Context, notifications []Notification) error
}
