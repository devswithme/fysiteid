import { logEvent } from "@/lib/utils";
import { authService } from "@/service/auth";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authService.logout(),
    onSuccess: () => {
      logEvent("auth", {
        type: "logout",
      });
      queryClient.clear();
      window.location.href = "/";
    },
    onError: (err) => {
      toast.error(err.message);
      window.location.href = "/";
    },
  });
}

export function useLogin() {
  logEvent("auth", {
    type: "login",
  });
  return {
    login: () => authService.login(),
  };
}
