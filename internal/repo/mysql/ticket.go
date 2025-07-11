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
	GetMany(ctx context.Context, accountId, offset, limit uint64) ([]*Ticket, error)
	GetOne(ctx context.Context, accountId uint64, ticketId uint64) (*Ticket, error)
	CreateOne(ctx context.Context, accountId uint64, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, ticketId uint64) (*Ticket, error)
}

type ticketDatabase struct {
	db *gorm.DB
}

func (t ticketDatabase) GetMany(ctx context.Context, accountId, offset, limit uint64) ([]*Ticket, error) {
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

func (t ticketDatabase) GetOne(ctx context.Context, accountId uint64, ticketId uint64) (*Ticket, error) {
	ticket := &Ticket{}

	res := t.db.Model(ticket).Where("id = ?", ticketId).Where("account_id = ?", accountId).Preload("Event").First(ticket)

	if res.Error != nil {
		return nil, res.Error
	}

	return ticket, nil
}

func (t ticketDatabase) CreateOne(ctx context.Context, accountId uint64, ticket *Ticket) (*Ticket, error) {
	res := t.db.Model(ticket).Create(ticket)

	if res.Error != nil {
		return nil, res.Error
	}

	return t.GetOne(ctx, accountId, ticket.ID)
}

func (t ticketDatabase) UpdateOne(ctx context.Context, ticketId uint64) (*Ticket, error) {
	ticket := &Ticket{}
	res := t.db.First(ticket, ticketId)
	if res.Error != nil {
		return nil, res.Error
	}

	ticket.Entered = true
	if err := t.db.Save(ticket).Error; err != nil {
		return nil, err
	}

	return ticket, nil
}

func NewTicketDatabase(db *gorm.DB) TicketDatabase {
	return &ticketDatabase{db: db}
}
