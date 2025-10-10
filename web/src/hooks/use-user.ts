import { logEvent } from "@/lib/utils";
import { userService } from "@/service/user";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

export const userKeys = {
  all: ["auth"] as const,
  user: () => [...userKeys.all, "user"] as const,
  public: (username: string) => [...userKeys.all, "public", username] as const,
  detail: (id: string) => [...userKeys.all, "detail", id] as const,
};

export function usePublicUser(username: string) {
  logEvent("user", {
    type: "getpublic",
  });
  return useQuery({
    queryKey: userKeys.public(username),
    queryFn: () => userService.getPublicUser(username),
    retry: false,
    staleTime: 5 * 60 * 1000,
    enabled: !!username,
  });
}

export function usePublicUserByTicketId(id: string) {
  logEvent("user", {
    type: "getpublicbyticket",
  });
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: () => userService.getUserByTicketId(id),
    retry: false,
    staleTime: 5 * 60 * 1000,
    enabled: !!id,
  });
}

export function useUser() {
  logEvent("user", {
    type: "get",
  });
  return useQuery({
    queryKey: userKeys.user(),
    queryFn: () => userService.getUser(),
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}

export function useUpdateUser() {
  logEvent("user", {
    type: "update",
  });
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: FormData) => userService.updateUser(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.user() });
    },
  });
}
