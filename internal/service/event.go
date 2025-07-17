package service

import (
	"Gober/internal/repo/mysql"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type CreateEventParams struct {
	Title    string    `json:"title"`
	Location string    `json:"location"`
	Date     time.Time `json:"date"` // ISO 8601 format
}

type EventService interface {
	GetEvents(ctx context.Context, offset uint64, limit uint64) ([]*mysql.Event, error)
	GetEvent(ctx context.Context, eventId uint64) (*mysql.Event, error)
	CreateEvent(ctx context.Context, event CreateEventParams) (*mysql.Event, error)
	UpdateEvent(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*mysql.Event, error)
	DeleteEvent(ctx context.Context, eventId uint64) error
}

type eventService struct {
	db mysql.EventDatabase
}

func (e eventService) GetEvents(ctx context.Context, offset uint64, limit uint64) ([]*mysql.Event, error) {
	events, err := e.db.GetEvents(ctx)

	if err != nil {
		return nil, status.Error(codes.DataLoss, "Không thể lấy danh sách sự kiện")
	}

	if offset >= uint64(len(events)) {
		return []*mysql.Event{}, status.Error(codes.InvalidArgument, "Số lượng sự kiện yêu cầu vượt quá tổng số sự kiện hiện có")
	}

	return events, nil
}

func (e eventService) GetEvent(ctx context.Context, eventId uint64) (*mysql.Event, error) {
	event, err := e.db.GetEvent(ctx, eventId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Không tìm thấy sự kiện với ID đã cho")
	}

	return event, nil
}

func (e eventService) CreateEvent(ctx context.Context, event CreateEventParams) (*mysql.Event, error) {
	newEvent := &mysql.Event{
		Title:    event.Title,
		Location: event.Location,
		Date:     event.Date,
	}

	createdEvent, err := e.db.CreateEvent(ctx, newEvent)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "Không thể tạo sự kiện mới")
	}

	return createdEvent, nil
}

func (e eventService) UpdateEvent(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*mysql.Event, error) {
	event, err := e.db.UpdateEvent(ctx, eventId, updateData)
	if err != nil {
		return nil, status.Error(codes.DataLoss, "Không thể cập nhật sự kiện")
	}

	if event == nil {
		return nil, status.Error(codes.NotFound, "Không tìm thấy sự kiện với ID đã cho")
	}

	return event, nil
}

func (e eventService) DeleteEvent(ctx context.Context, eventId uint64) error {
	err := e.db.DeleteEvent(ctx, eventId)
	if err != nil {
		return status.Error(codes.DataLoss, "Không thể xóa sự kiện")
	}

	return nil
}

func NewEventService(db mysql.EventDatabase) EventService {
	return &eventService{
		db: db,
	}
}
