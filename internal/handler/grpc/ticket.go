package grpc

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/service"
	"context"
)

type TicketHandler struct {
	gober.UnimplementedGoberServiceServer
	TicketService service.TicketService
}

func (h *TicketHandler) GetMany(ctx context.Context, req *gober.ListTicketsRequest) (*gober.ListTicketsResponse, error) {
	tickets, err := h.TicketService.GetMany(ctx, req.AccountId, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	var ticketResponses []*gober.Ticket
	for _, ticket := range tickets {
		ticketResponses = append(ticketResponses, &gober.Ticket{
			TicketId:  ticket.ID,
			EventId:   ticket.EventID,
			AccountId: ticket.AccountID,
			Entered:   ticket.Entered,
			CreatedAt: ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: ticket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &gober.ListTicketsResponse{
		Tickets: ticketResponses,
		TotalTicketCount: uint64(len(tickets)),
		HasNext: len(tickets) > int(req.Offset+req.Limit),
	}, nil
}

func (h *TicketHandler) GetOne(ctx context.Context, req *gober.GetTicketRequest) (*gober.GetTicketResponse, error) {
	ticket, err := h.TicketService.GetOne(ctx, service.TicketParams{
		TicketId:  req.TicketId,
		AccountId: req.AccountId,
	})
	if err != nil {
		return nil, err
	}

	return &gober.GetTicketResponse{
		Ticket: &gober.Ticket{
			TicketId:  ticket.ID,
			EventId:   ticket.EventID,
			AccountId: ticket.AccountID,
			Entered:   ticket.Entered,
			CreatedAt: ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: ticket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (h *TicketHandler) CreateOne(ctx context.Context, req *gober.CreateTicketRequest) (*gober.CreateTicketResponse, error) {
	ticket, err := h.TicketService.CreateOne(ctx, service.CreateTicketParams{
		EventId:   req.EventId,
		AccountId: req.AccountId,
	})
	if err != nil {
		return nil, err
	}

	return &gober.CreateTicketResponse{
		Ticket: &gober.Ticket{
			TicketId:  ticket.ID,
			EventId:   ticket.EventID,
			AccountId: ticket.AccountID,
			Entered:   ticket.Entered,
			CreatedAt: ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: ticket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (h *TicketHandler) UpdateOne(ctx context.Context, req *gober.UpdateTicketRequest) (*gober.UpdateTicketResponse, error) {
	ticket, err := h.TicketService.UpdateOne(ctx, service.TicketParams{
		TicketId:  req.TicketId,
		AccountId: req.AccountId,
	})
	if err != nil {
		return nil, err
	}

	return &gober.UpdateTicketResponse{
		Ticket: &gober.Ticket{
			TicketId:  ticket.ID,
			EventId:   ticket.EventID,
			AccountId: ticket.AccountID,
			Entered: ticket.Entered,
			CreatedAt: ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: ticket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}
