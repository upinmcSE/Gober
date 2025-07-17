package http

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/utils/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type EventHandler struct {
	client gober.GoberServiceClient
}

func (eh *EventHandler) CreateEventHandler(c *gin.Context) {
	var req gober.EventUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid request",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}
	request := gober.CreateEventRequest{
		EventUpdate: &req,
	}

	resp, err := eh.client.CreateEvent(c, &request)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 201, "Event created successfully", resp)
}

func (eh *EventHandler) GetEventHandler(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Event ID is required",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	id, err := strconv.ParseUint(eventID, 10, 64)

	resp, err := eh.client.GetEvent(c, &gober.GetEventRequest{EventId: id})
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Event retrieved successfully", resp)
}

func (eh *EventHandler) ListEventsHandler(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")

	offset, err1 := strconv.ParseUint(offsetStr, 10, 64)
	limit, err2 := strconv.ParseUint(limitStr, 10, 64)

	if err1 != nil || err2 != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid offset or limit",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	// 2. Gán vào proto struct
	req := &gober.ListEventsRequest{
		Offset: offset,
		Limit:  limit,
	}

	resp, err := eh.client.GetEvents(c, req)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Events listed successfully", resp)
}

func (eh *EventHandler) UpdateEventHandler(c *gin.Context) {
	var req gober.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid request",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Event ID is required",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	id, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid Event ID format",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	req.EventId = id

	resp, err := eh.client.UpdateEvent(c, &req)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Event updated successfully", resp)
}

func (eh *EventHandler) DeleteEventHandler(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Event ID is required",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	id, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid Event ID format",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	resp, err := eh.client.DeleteEvent(c, &gober.DeleteEventRequest{EventId: id})
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Event deleted successfully", resp)
}

func NewEventHandler(client gober.GoberServiceClient) *EventHandler {
	return &EventHandler{
		client: client,
	}
}
