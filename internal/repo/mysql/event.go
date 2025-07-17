package mysql

import (
	"context"
	"fmt"
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
	GetEvents(ctx context.Context) ([]*Event, error)
	GetEvent(ctx context.Context, eventId uint64) (*Event, error)
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	UpdateEvent(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*Event, error)
	DeleteEvent(ctx context.Context, eventId uint64) error
}

type eventDatabase struct {
	db *gorm.DB
}

func (e eventDatabase) GetEvents(ctx context.Context) ([]*Event, error) {
	var events []*Event

	res := e.db.Model(&Event{}).Order("created_at desc").Find(&events)

	if res.Error != nil {
		return nil, res.Error
	}

	fmt.Println(events)

	if res.RowsAffected == 0 {
		return []*Event{}, nil
	}
	return events, nil
}

func (e eventDatabase) GetEvent(ctx context.Context, eventId uint64) (*Event, error) {
	event := &Event{}

	res := e.db.Model(event).Where("id = ?", eventId).First(event)

	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (e eventDatabase) CreateEvent(ctx context.Context, event *Event) (*Event, error) {
	res := e.db.Model(event).Create(event)

	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (e eventDatabase) UpdateEvent(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*Event, error) {
	event := new(Event)

	updateRes := e.db.WithContext(ctx).Where("id = ?", eventId).Updates(updateData)

	if updateRes.Error != nil {
		return nil, updateRes.Error
	}

	getRes := e.db.Model(event).Where("id = ?", eventId).First(event)

	if getRes.Error != nil {
		return nil, getRes.Error
	}

	return event, nil
}

func (e eventDatabase) DeleteEvent(ctx context.Context, eventId uint64) error {
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
