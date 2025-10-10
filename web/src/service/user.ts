/* eslint-disable @typescript-eslint/no-explicit-any */
import type { User } from "@/types/user";
import { type APIResponse } from "@/types/api";
import { fetch } from "../lib/fetch";

export class UserService {
  async getUser(): Promise<User> {
    try {
      const { data } = await fetch.get<APIResponse<User>>("/user/me");
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async getUserByTicketId(id: string): Promise<User> {
    try {
      const { data } = await fetch.get<APIResponse<User>>(
        `/public/user/id/${id}`
      );
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async updateUser(payload: FormData): Promise<null> {
    try {
      const { data } = await fetch.patch<APIResponse<null>>(
        "/user/me",
        payload
      );
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }

  async getPublicUser(username: string): Promise<User> {
    try {
      const { data } = await fetch.get<APIResponse<User>>(
        `/public/user/${username}`
      );
      return data.data;
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }
}

export const userService = new UserService();
