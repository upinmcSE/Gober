package grpc

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/service"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type GoberHandler struct {
	gober.UnimplementedGoberServiceServer
	AccountService service.AccountService
	EventService   service.EventService
	TicketService  service.TicketService
}

func (h *GoberHandler) CreateAccount(ctx context.Context, req *gober.CreateAccountRequest) (*gober.CreateAccountResponse, error) {
	params := service.CreateAccountParams{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.AccountService.CreateAccount(ctx, params)
	if err != nil {
		return nil, err
	}

	return &gober.CreateAccountResponse{
		AccountId: output.ID,
	}, nil
}

func (h *GoberHandler) CreateSession(ctx context.Context, req *gober.CreateSessionRequest) (*gober.CreateSessionResponse, error) {
	params := service.CreateSessionParams{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.AccountService.CreateSession(ctx, params)
	if err != nil {
		return nil, err
	}
	fmt.Println("ouput:", output.Account)
	fmt.Println("ouput:", output.AccessToken)
	fmt.Println("ouput:", output.RefreshToken)

	return &gober.CreateSessionResponse{
		OfAccount:    output.Account,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (h *GoberHandler) RefreshSession(ctx context.Context, req *gober.RefreshSessionRequest) (*gober.RefreshSessionResponse, error) {
	output, err := h.AccountService.RefreshSession(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &gober.RefreshSessionResponse{
		OfAccount:    output.Account,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (h *GoberHandler) GetAccount(ctx context.Context, req *gober.GetAccountRequest) (*gober.GetAccountResponse, error) {
	account, err := h.AccountService.GetAccountByID(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &gober.GetAccountResponse{
		OfAccount: account,
	}, nil
}

func (h *GoberHandler) DeleteSession(ctx context.Context, req *gober.DeleteSessionRequest) (*gober.DeleteSessionResponse, error) {
	if err := h.AccountService.DeleteSession(ctx, req.AccessToken, req.RefreshToken); err != nil {
		return nil, err
	}

	return &gober.DeleteSessionResponse{}, nil
}

func (h *GoberHandler) GetEvents(ctx context.Context, req *gober.ListEventsRequest) (*gober.ListEventsResponse, error) {
	events, err := h.EventService.GetEvents(ctx, req.Offset, req.Limit)
	fmt.Println("handler grpc", events)
	if err != nil {
		return nil, err
	}

	var eventResponses []*gober.Event
	for _, event := range events {
		eventResponses = append(eventResponses, &gober.Event{
			EventId:               event.ID,
			Title:                 event.Title,
			Location:              event.Location,
			TotalTicketsPurchased: event.TotalTicketsPurchased,
			TotalTicketsEntered:   event.TotalTicketsEntered,
			Date:                  event.Date.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:             event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:             event.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	fmt.Println(req.Offset + req.Limit)

	return &gober.ListEventsResponse{
		Events:          eventResponses,
		TotalEventCount: uint64(len(events)),
		HasNext:         len(events) > int(req.Offset+req.Limit),
	}, nil
}

func (h *GoberHandler) GetEvent(ctx context.Context, req *gober.GetEventRequest) (*gober.GetEventResponse, error) {
	event, err := h.EventService.GetEvent(ctx, req.EventId)
	if err != nil {
		return nil, err
	}

	return &gober.GetEventResponse{
		Event: &gober.Event{
			EventId:               event.ID,
			Title:                 event.Title,
			Location:              event.Location,
			TotalTicketsPurchased: event.TotalTicketsPurchased,
			TotalTicketsEntered:   event.TotalTicketsEntered,
			Date:                  event.Date.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:             event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:             event.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (h *GoberHandler) CreateEvent(ctx context.Context, req *gober.CreateEventRequest) (*gober.CreateEventResponse, error) {
	date, err := time.Parse(time.RFC3339, req.EventUpdate.Date)

	if err != nil {
		// Try to parse just the date (YYYY-MM-DD)
		parsedDate, err2 := time.Parse("2006-01-02", req.EventUpdate.Date)
		if err2 != nil {
			return nil, status.Error(codes.InvalidArgument, "Date must be in RFC3339 format or YYYY-MM-DD")
		}
		// Add default time 20:00:00 (8 PM)
		date = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 20, 0, 0, 0, time.UTC)
	}

	event, err := h.EventService.CreateEvent(ctx, service.CreateEventParams{
		Title:    req.EventUpdate.Title,
		Location: req.EventUpdate.Location,
		Date:     date,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create event: "+err.Error())
	}

	return &gober.CreateEventResponse{
		Event: &gober.Event{
			EventId:               event.ID,
			Title:                 event.Title,
			Location:              event.Location,
			TotalTicketsPurchased: event.TotalTicketsPurchased,
			TotalTicketsEntered:   event.TotalTicketsEntered,
			Date:                  event.Date.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:             event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:             event.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (h *GoberHandler) UpdateEvent(ctx context.Context, req *gober.UpdateEventRequest) (*gober.UpdateEventResponse, error) {
	updateData := make(map[string]interface{})

	if req.EventUpdate.Title != "" {
		updateData["title"] = req.EventUpdate.Title
	}

	if req.EventUpdate.Location != "" {
		updateData["location"] = req.EventUpdate.Location
	}

	if req.EventUpdate.Date != "" {
		date, err := time.Parse(time.RFC3339, req.EventUpdate.Date)
		if err != nil {
			// Try to parse just the date (YYYY-MM-DD)
			parsedDate, err2 := time.Parse("2006-01-02", req.EventUpdate.Date)
			if err2 != nil {
				return nil, status.Error(codes.InvalidArgument, "Date must be in RFC3339 format or YYYY-MM-DD")
			}
			// Add default time 20:00:00 (8 PM)
			date = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 20, 0, 0, 0, time.UTC)
		}
		updateData["date"] = date
	}

	event, err := h.EventService.UpdateEvent(ctx, req.EventId, updateData)
	if err != nil {
		return nil, err
	}

	return &gober.UpdateEventResponse{
		Event: &gober.Event{
			EventId:               event.ID,
			Title:                 event.Title,
			Location:              event.Location,
			TotalTicketsPurchased: event.TotalTicketsPurchased,
			TotalTicketsEntered:   event.TotalTicketsEntered,
			Date:                  event.Date.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:             event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:             event.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (h *GoberHandler) DeleteEvent(ctx context.Context, req *gober.DeleteEventRequest) (*gober.DeleteEventResponse, error) {
	if err := h.EventService.DeleteEvent(ctx, req.EventId); err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete event: "+err.Error())
	}

	return &gober.DeleteEventResponse{}, nil
}

// Ticket Methods
func (h *GoberHandler) GetTickets(ctx context.Context, req *gober.ListTicketsRequest) (*gober.ListTicketsResponse, error) {
	tickets, err := h.TicketService.GetTickets(ctx, req.AccountId, req.Offset, req.Limit)
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
		Tickets:          ticketResponses,
		TotalTicketCount: uint64(len(tickets)),
		HasNext:          len(tickets) > int(req.Offset+req.Limit),
	}, nil
}

func (h *GoberHandler) GetTicket(ctx context.Context, req *gober.GetTicketRequest) (*gober.GetTicketResponse, error) {
	ticket, err := h.TicketService.GetTicket(ctx, service.TicketParams{
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

func (h *GoberHandler) CreateTicket(ctx context.Context, req *gober.CreateTicketRequest) (*gober.CreateTicketResponse, error) {
	ticket, err := h.TicketService.CreateTicket(ctx, service.CreateTicketParams{
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

func (h *GoberHandler) UpdateTicket(ctx context.Context, req *gober.UpdateTicketRequest) (*gober.UpdateTicketResponse, error) {
	ticket, err := h.TicketService.UpdateTicket(ctx, service.TicketParams{
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
			Entered:   ticket.Entered,
			CreatedAt: ticket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: ticket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func NewGoberHandler(accountService service.AccountService, eventService service.EventService, ticketService service.TicketService) (gober.GoberServiceServer, error) {
	return &GoberHandler{
		AccountService: accountService,
		EventService:   eventService,
		TicketService:  ticketService,
	}, nil
}
