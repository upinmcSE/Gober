import { ApiResponse } from "./api";
import { Event } from "./event";

export type TicketResponse = ApiResponse<Ticket>;
export type TicketListResponse = ApiResponse<TicketListData>;

export type Ticket = {
  ticket_id: number;
  event_id: number;
  account_id: number;
  event: Event;
  entered: boolean;
  createdAt: string;
  updatedAt: string;
}

export type TicketListData = {
  tickets: Ticket[];
  total_event_count: number;
};