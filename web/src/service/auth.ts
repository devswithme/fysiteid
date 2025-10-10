/* eslint-disable @typescript-eslint/no-explicit-any */
import { fetch } from "../lib/fetch";

export class AuthService {
  async login(path?: string): Promise<void> {
    window.location.href = `${
      import.meta.env.VITE_API_URL
    }/login/google?redirect=${path}`;
  }

  async logout(): Promise<void> {
    try {
      await fetch.post("/auth/logout");
    } catch (err: any) {
      throw new Error(err?.response?.data?.meta?.message);
    }
  }
}

export const authService = new AuthService();
