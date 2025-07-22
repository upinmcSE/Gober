package http

import (
	"Gober/internal/generated/grpc/gober"
	"Gober/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TicketHandler struct {
	client gober.GoberServiceClient
}

func (th *TicketHandler) CreateTicketHandler(c *gin.Context) {
	var req gober.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid request",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	value, _ := c.Get("accountID")
	accountID, _ := value.(uint64)
	req.AccountId = accountID

	fmt.Println("http:", accountID, req.EventId)

	resp, err := th.client.CreateTicket(c, &req)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 201, "Ticket created successfully", resp)
}

func (th *TicketHandler) GetTicketHandler(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Ticket ID is required",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	id, err := strconv.ParseUint(ticketID, 10, 64)
	if err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid Ticket ID format",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	value, _ := c.Get("accountID")
	accountID, _ := value.(uint64)

	resp, err := th.client.GetTicket(c, &gober.GetTicketRequest{
		TicketId:  id,
		AccountId: accountID,
	})
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Ticket retrieved successfully", resp)
}

func (th *TicketHandler) ListTicketsHandler(c *gin.Context) {
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

	value, _ := c.Get("accountID")
	accountID, _ := value.(uint64)

	req := &gober.ListTicketsRequest{
		AccountId: accountID,
		Offset:    offset,
		Limit:     limit,
	}

	resp, err := th.client.GetTickets(c, req)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Tickets retrieved successfully", resp)
}

func (th *TicketHandler) UpdateTicketHandler(c *gin.Context) {
	req := &gober.UpdateTicketRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ResponseError(c, 400, &response.AppError{
			Message: "Invalid request",
			Code:    response.ErrCodeBadRequest,
		})
		return
	}

	resp, err := th.client.UpdateTicket(c, req)
	if err != nil {
		response.HandleGrpcError(c, err)
		return
	}

	response.ResponseSuccess(c, 200, "Ticket updated successfully", resp)
}

func NewTicketHandler(client gober.GoberServiceClient) *TicketHandler {
	return &TicketHandler{
		client: client,
	}
}
