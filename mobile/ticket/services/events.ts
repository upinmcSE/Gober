import { EventListResponse, EventResponse } from "@/types/event";
import { Api } from "./api";

async function createOne(title: string, location: string, date: string): Promise<EventResponse> {
  return Api.post("/events/create", { title, location, date });
}

async function getOne(id: number): Promise<EventResponse> {
  return Api.get(`/events/${id}`);
}

async function getAll(): Promise<EventListResponse> {
  return Api.get("/events?offset=0&limit=10");
}

async function updateOne(id: number, title: string, location: string, date: string): Promise<EventResponse> {
  return Api.put(`/events/${id}`, { title, location, date });
} 

async function deleteOne(id: number): Promise<EventResponse> {
  return Api.delete(`/events/${id}`);
}

const eventService = {
  createOne,
  getOne,
  getAll,
  updateOne,
  deleteOne,
};

export { eventService };

