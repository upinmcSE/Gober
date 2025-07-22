package mysql

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Ticket struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	EventID   uint64    `json:"eventId"`
	AccountID uint64    `json:"userId" gorm:"foreignkey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Event     Event     `json:"event" gorm:"foreignkey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entered   bool      `json:"entered" default:"false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TicketDatabase interface {
	GetTickets(ctx context.Context, accountId, offset, limit uint64) ([]*Ticket, error)
	GetTicket(ctx context.Context, accountId uint64, ticketId uint64) (*Ticket, error)
	CreateTicket(ctx context.Context, accountId uint64, ticket *Ticket) (*Ticket, error)
	UpdateTicket(ctx context.Context, ticketId uint64) (*Ticket, error)
}

type ticketDatabase struct {
	db *gorm.DB
}

func (t ticketDatabase) GetTickets(ctx context.Context, accountId, offset, limit uint64) ([]*Ticket, error) {
	var tickets []*Ticket

	res := t.db.
		Model(&Ticket{}).
		Where("account_id = ?", accountId).
		Preload("Event").
		Order("updated_at desc").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&tickets)

	if res.Error != nil {
		return nil, res.Error
	}

	return tickets, nil
}

func (t ticketDatabase) GetTicket(ctx context.Context, accountId uint64, ticketId uint64) (*Ticket, error) {
	ticket := &Ticket{}

	res := t.db.Model(ticket).Where("id = ?", ticketId).Where("account_id = ?", accountId).Preload("Event").First(ticket)

	if res.Error != nil {
		return nil, res.Error
	}

	return ticket, nil
}

func (t ticketDatabase) CreateTicket(ctx context.Context, accountId uint64, ticket *Ticket) (*Ticket, error) {
	err := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Tạo ticket
		if err := tx.Create(ticket).Error; err != nil {
			return err
		}

		// Cập nhật total_tickets_purchased trong Event
		if err := tx.Model(&Event{}).
			Where("id = ?", ticket.EventID).
			Update("total_tickets_purchased", gorm.Expr("total_tickets_purchased + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	result, err := t.GetTicket(ctx, accountId, ticket.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t ticketDatabase) UpdateTicket(ctx context.Context, ticketId uint64) (*Ticket, error) {
	var updatedTicket *Ticket

	err := t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ticket := &Ticket{}
		if err := tx.First(ticket, ticketId).Error; err != nil {
			return err
		}

		ticket.Entered = true
		if err := tx.Save(ticket).Error; err != nil {
			return err
		}

		if err := tx.Model(&Event{}).
			Where("id = ?", ticket.EventID).
			Update("total_tickets_entered", gorm.Expr("total_tickets_entered + ?", 1)).Error; err != nil {
			return err
		}

		updatedTicket = ticket
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedTicket, nil
}

func NewTicketDatabase(db *gorm.DB) TicketDatabase {
	return &ticketDatabase{db: db}
}
