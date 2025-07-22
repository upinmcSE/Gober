import { ApiResponse } from "@/types/api";
import { Ticket, TicketListData, TicketResponse } from "@/types/ticket";
import { Api } from "./api";

async function createOne(eventId: number): Promise<TicketResponse> {
  return Api.post("/tickets", { eventId });
}

async function getOne(id: number): Promise<ApiResponse<{ticket: Ticket, qrcode: string}>> {
  return Api.get(`/tickets/${id}`);
}

async function getAll(): Promise<ApiResponse<TicketListData>> {
  return Api.get("/tickets?offset=0&limit=10");
}

async function validateOne(ticketId: number, ownerId: number): Promise<ApiResponse<Ticket>> {
  return Api.post("/tickets/validate", { account_id: ownerId, ticket_id: ticketId });
}

const ticketService = {
  createOne,
  getOne,
  getAll,
  validateOne,
};

export { ticketService };

