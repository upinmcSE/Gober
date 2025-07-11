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
	GetMany(ctx context.Context, offset uint64, limit uint64) ([]*mysql.Event, error)
	GetOne(ctx context.Context, eventId uint64) (*mysql.Event, error)
	CreateOne(ctx context.Context, event CreateEventParams) (*mysql.Event, error)
	UpdateOne(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*mysql.Event, error)
	DeleteOne(ctx context.Context, eventId uint64) error
}

type eventService struct {
	db mysql.EventDatabase
}

func (e eventService) GetMany(ctx context.Context, offset uint64, limit uint64) ([]*mysql.Event, error) {
	events, err := e.db.GetMany(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Không thể lấy danh sách sự kiện")
	}

	if offset >= uint64(len(events)) {
		return []*mysql.Event{}, status.Error(codes.InvalidArgument, "Số lượng sự kiện yêu cầu vượt quá tổng số sự kiện hiện có")
	}

	end := offset + limit
	if end > uint64(len(events)) {
		end = uint64(len(events))
	}

	return events[offset:end], nil
}

func (e eventService) GetOne(ctx context.Context, eventId uint64) (*mysql.Event, error) {
	event, err := e.db.GetOne(ctx, eventId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Không tìm thấy sự kiện với ID đã cho")
	}

	return event, nil
}

func (e eventService) CreateOne(ctx context.Context, event CreateEventParams) (*mysql.Event, error) {
	newEvent := &mysql.Event{
		Title:    event.Title,
		Location: event.Location,
		Date:     event.Date,
	}

	createdEvent, err := e.db.CreateOne(ctx, newEvent)
	if err != nil {
		return nil, status.Error(codes.Internal, "Không thể tạo sự kiện mới")
	}

	return createdEvent, nil
}

func (e eventService) UpdateOne(ctx context.Context, eventId uint64, updateData map[string]interface{}) (*mysql.Event, error) {
	event, err := e.db.UpdateOne(ctx, eventId, updateData)
	if err != nil {
		return nil, status.Error(codes.Internal, "Không thể cập nhật sự kiện")
	}

	if event == nil {
		return nil, status.Error(codes.NotFound, "Không tìm thấy sự kiện với ID đã cho")
	}

	return event, nil
}

func (e eventService) DeleteOne(ctx context.Context, eventId uint64) error {
	err := e.db.DeleteOne(ctx, eventId)
	if err != nil {
		return status.Error(codes.Internal, "Không thể xóa sự kiện")
	}

	return nil
}

func NewEventService(db mysql.EventDatabase) EventService {
	return &eventService{
		db: db,
	}
}
