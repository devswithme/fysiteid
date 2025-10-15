import { userService } from "@/service/user";
import { redirect } from "@tanstack/react-router";

export const authMiddleware = {
  requireAuth: async () => {
    try {
      await userService.getUser();
    } catch {
      throw redirect({
        to: "/",
        search: {
          redirect: window.location.pathname,
        },
      });
    }
  },
};
