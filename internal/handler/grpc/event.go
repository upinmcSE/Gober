package grpc

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/internal/service"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type EventHandler struct {
	gober.UnimplementedGoberServiceServer
	EventService service.EventService
}

func (h *EventHandler) GetMany(ctx context.Context, req *gober.ListEventsRequest) (*gober.ListEventsResponse, error) {
	events, err := h.EventService.GetMany(ctx, req.Offset, req.Limit)
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

	return &gober.ListEventsResponse{
		Events:          eventResponses,
		TotalEventCount: uint64(len(events)),
		HasNext:         len(events) > int(req.Offset+req.Limit),
	}, nil
}

func (h *EventHandler) GetOne(ctx context.Context, req *gober.GetEventRequest) (*gober.GetEventResponse, error) {
	event, err := h.EventService.GetOne(ctx, req.EventId)
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

func (h *EventHandler) CreateOne(ctx context.Context, req *gober.CreateEventRequest) (*gober.CreateEventResponse, error) {
	date, err := time.Parse(time.RFC3339, req.EventUpdate.Date)
	if err != nil {
		return nil, err
	}
	event, err := h.EventService.CreateOne(ctx, service.CreateEventParams{
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

func (h *EventHandler) UpdateOne(ctx context.Context, req *gober.UpdateEventRequest) (*gober.UpdateEventResponse, error) {
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
			return nil, status.Error(codes.InvalidArgument, "Invalid date format")
		}
		updateData["date"] = date
	}

	event, err := h.EventService.UpdateOne(ctx, req.EventId, updateData)
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

func (h *EventHandler) DeleteOne(ctx context.Context, req *gober.DeleteEventRequest) (*gober.DeleteEventResponse, error) {
	if err := h.EventService.DeleteOne(ctx, req.EventId); err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete event: "+err.Error())
	}

	return &gober.DeleteEventResponse{}, nil
}

func NewEventHandler(eventService service.EventService) (gober.GoberServiceServer, error) {
	return &EventHandler{
		EventService: eventService,
	}, nil
}
