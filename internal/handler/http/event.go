package http

import "Gober/internal/generated/grpc/gober"

type EventHandler struct {
	client gober.GoberServiceClient
}

func NewEventHandler(client gober.GoberServiceClient) *EventHandler {
	return &EventHandler{
		client: client,
	}
}
