package models

import (
	"context"
	"time"
)

type Event struct {
	ID                    uint `json:"id"`
	Name                  string `json:"name"`
	Location              string `json:"location"`
	TotalTicketsPurchased int    `json:"total_tickets_purchased"`
	TotalTicketsEntered   int    `json:"total_tickets_entered"`
	Date                  time.Time `json:"date"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type EventRepository interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventId uint) (*Event, error)
	CreateOne(ctx context.Context, event *Event) (*Event, error)
	UpdateOne(ctx context.Context, event *Event) (*Event, error)
	DeleteOne(ctx context.Context, eventId uint) error
}