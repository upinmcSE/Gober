package models

import "context"

type Ticket struct {
	ID        string `json:"id"`
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	Event     Event  `json:"event"`
	Entered   bool   `json:"entered"`
	CreatedAt string `json:"created_at"`
}

type TicketRepository interface {
	GetMany(ctx context.Context, userId uint) ([]*Ticket, error)
	GetOne(ctx context.Context, userId uint, ticketId uint) (*Ticket, error)
	CreateOne(ctx context.Context, userId uint, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, userId uint, ticketId uint, updateData map[string]interface{}) (*Ticket, error)
}

type ValidateTicket struct {
	TicketId uint `json:"ticketId"`
	OwnerId  uint `json:"ownerId"`
}