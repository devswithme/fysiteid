/* eslint-disable @typescript-eslint/no-explicit-any */
import { fetch } from "@/lib/fetch";
import type { APIResponse } from "@/types/api";
import type {
  PaginatedRegistrants,
  RegistrantGetByUserID,
} from "@/types/registrant";

export class RegistrantService {
  async createRegistrant(ticketID: string, state: string) {
    try {
      const { data } = await fetch.post<APIResponse<null>>(
        `/registrant/${ticketID}?state=${state}`
      );
      return data.meta.message;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async getRegistrantByTicketID(
    ticketID: string,
    page: number = 1,
    limit: number = 10,
    search: string
  ) {
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
        search,
      });

      const { data } = await fetch.get<APIResponse<PaginatedRegistrants>>(
        `/registrant/${ticketID}?${params.toString()}`
      );

      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async generateState(ticketID: string) {
    try {
      const { data } = await fetch.get<APIResponse<string>>(
        `/registrant/gen/${ticketID}`
      );
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async getRegistrantByUserID() {
    try {
      const { data } = await fetch.get<APIResponse<RegistrantGetByUserID[]>>(
        "/registrant"
      );
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }
}

export const registrantService = new RegistrantService();
