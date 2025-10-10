/* eslint-disable @typescript-eslint/no-explicit-any */
import { fetch } from "@/lib/fetch";
import type { APIResponse } from "@/types/api";
import type { Ticket } from "@/types/ticket";

export class TicketService {
  async getTickets(): Promise<Ticket[]> {
    try {
      const { data } = await fetch.get<APIResponse<Ticket[]>>("/ticket");
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async getPublicTickets(username: string): Promise<Ticket[]> {
    try {
      const { data } = await fetch.get<APIResponse<Ticket[]>>(
        `/public/ticket/${username}`
      );
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async getPublicTicket(id: string): Promise<Ticket> {
    try {
      const { data } = await fetch.get<APIResponse<Ticket>>(
        `/public/ticket/id/${id}`
      );
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async getTicketById(id: string): Promise<Ticket> {
    try {
      const { data } = await fetch.get<APIResponse<Ticket>>(`/ticket/${id}`);
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async createTicket(payload: FormData): Promise<string> {
    try {
      const { data } = await fetch.post<APIResponse<null>>("/ticket", payload);
      return data.meta.message;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async updateTicket(id: string, payload: FormData): Promise<string> {
    try {
      const { data } = await fetch.patch<APIResponse<null>>(
        `/ticket/${id}`,
        payload
      );
      return data.meta.message;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async deleteTicket(id: string): Promise<string> {
    try {
      const { data } = await fetch.delete<APIResponse<null>>(`/ticket/${id}`);
      return data.meta.message;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }
}

export const ticketService = new TicketService();
