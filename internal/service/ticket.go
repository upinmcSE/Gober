package service

import (
	"Gober/internal/repo/mysql"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type CreateTicketParams struct {
	EventId   uint64 `json:"eventId"`
	AccountId uint64 `json:"accountId"`
}

type TicketParams struct {
	TicketId  uint64 `json:"ticketId"`
	AccountId uint64 `json:"accountId"`
}

type TicketService interface {
	GetMany(ctx context.Context, accountId uint64, offset uint64, limit uint64) ([]*mysql.Ticket, error)
	GetOne(ctx context.Context, params TicketParams) (*mysql.Ticket, error)
	CreateOne(ctx context.Context, params CreateTicketParams) (*mysql.Ticket, error)
	UpdateOne(ctx context.Context, params TicketParams) (*mysql.Ticket, error)
}

type ticketService struct {
	db mysql.TicketDatabase
}

func (t ticketService) GetMany(ctx context.Context, accountId uint64, offset uint64, limit uint64) ([]*mysql.Ticket, error) {
	tickets, err := t.db.GetMany(ctx, accountId, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "Không thể lấy danh sách vé")
	}

	return tickets, nil
}

func (t ticketService) GetOne(ctx context.Context, params TicketParams) (*mysql.Ticket, error) {
	ticket, err := t.db.GetOne(ctx, params.AccountId, params.TicketId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy vé")
		}
		return nil, status.Error(codes.Internal, "Lỗi hệ thống")
	}

	return ticket, nil
}

func (t ticketService) CreateOne(ctx context.Context, params CreateTicketParams) (*mysql.Ticket, error) {
	ticket := &mysql.Ticket{
		EventID:   params.EventId,
		AccountID: params.AccountId,
	}

	createdTicket, err := t.db.CreateOne(ctx, params.AccountId, ticket)
	if err != nil {
		return nil, status.Error(codes.Internal, "Không thể tạo vé mới")
	}

	return createdTicket, nil
}

func (t ticketService) UpdateOne(ctx context.Context, params TicketParams) (*mysql.Ticket, error) {
	ticket, err := t.db.GetOne(ctx, params.AccountId, params.TicketId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy vé với ID")
		}
		return nil, status.Error(codes.Internal, "Lỗi hệ thống")
	}

	if ticket == nil {
		return nil, status.Error(codes.NotFound, "Không tìm thấy vé với ID đã cho")
	}

	if !ticket.Entered {
		ticket, err = t.db.UpdateOne(ctx, params.TicketId)
		if err != nil {
			return nil, status.Error(codes.Internal, "Không thể cập nhật vé")
		}
	}
	return ticket, nil
}

func NewTicketService(db mysql.TicketDatabase) TicketService {
	return &ticketService{db: db}
}
