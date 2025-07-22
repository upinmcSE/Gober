import { ApiResponse } from "./api";

export type EventResponse = ApiResponse<Event>;
export type EventListResponse = ApiResponse<EventListData>;

export type Event = {
  event_id: number;
  title: string;
  location: string;
  total_tickets_purchased: number;
  total_tickets_entered: number;
  date: string;
  created_at: string;
  updated_at: string;
}

export type EventListData = {
  events: Event[];
  total_event_count: number;
};