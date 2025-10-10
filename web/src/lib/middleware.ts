import { userService } from "@/service/user";
import { redirect } from "@tanstack/react-router";

export const authMiddleware = {
  requireAuth: async () => {
    if (!userService.getUser()) {
      const user = await userService.getUser();

      if (!user) {
        throw redirect({
          to: "/",
          search: {
            redirect: window.location.pathname,
          },
        });
      }
    }
  },
};
