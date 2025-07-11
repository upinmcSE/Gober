package mysql

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	ID                    uint64    `json:"id" gorm:"primarykey"`
	Title                 string    `json:"title"`
	Location              string    `json:"location"`
	TotalTicketsPurchased uint64    `json:"totalTicketsPurchased" gorm:"-"`
	TotalTicketsEntered   uint64    `json:"totalTicketsEntered" gorm:"-"`
	Date                  time.Time `json:"date"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type EventDatabase interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventId uint64) (*Event, error)
	CreateOne(ctx context.Context, event *Event) (*Event, error)
	UpdateOne(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*Event, error)
	DeleteOne(ctx context.Context, eventId uint64) error
}

type eventDatabase struct {
	db *gorm.DB
}

func (e eventDatabase) GetMany(ctx context.Context) ([]*Event, error) {
	var events []*Event

	res := e.db.Model(&Event{}).Order("updated_at desc").Find(&events)

	if res.Error != nil {
		return nil, res.Error
	}

	return events, nil
}

func (e eventDatabase) GetOne(ctx context.Context, eventId uint64) (*Event, error) {
	event := &Event{}

	res := e.db.Model(event).Where("id = ?", eventId).First(event)

	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (e eventDatabase) CreateOne(ctx context.Context, event *Event) (*Event, error) {
	res := e.db.Model(event).Create(event)

	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (e eventDatabase) UpdateOne(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*Event, error) {
	event := &Event{}

	updateRes := e.db.Model(event).Where("id = ?", eventId).Updates(updateData)

	if updateRes.Error != nil {
		return nil, updateRes.Error
	}

	getRes := e.db.Model(event).Where("id = ?", eventId).First(event)

	if getRes.Error != nil {
		return nil, getRes.Error
	}

	return event, nil
}

func (e eventDatabase) DeleteOne(ctx context.Context, eventId uint64) error {
	res := e.db.Delete(&Event{}, eventId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (e *Event) AfterFind(db *gorm.DB) (err error) {
	//baseQuery := db.Model(&Ticket{}).Where(&Ticket{EventID: e.ID})
	//
	//if res := baseQuery.Count(&e.TotalTicketsPurchased); res.Error != nil {
	//	return res.Error
	//}
	//if res := baseQuery.Where("entered = ?", true).Count(&e.TotalTicketsEntered); res.Error != nil {
	//	return res.Error
	//}
	return nil
}

func NewEventDatabase(db *gorm.DB) EventDatabase {
	return &eventDatabase{
		db: db,
	}
}
