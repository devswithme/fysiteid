import { logEvent } from "@/lib/utils";
import { ticketService } from "@/service/ticket";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export const ticketKeys = {
  all: ["ticket"] as const,
  list: () => [...ticketKeys.all, "list"] as const,
  listByUsername: (username: string) =>
    [...ticketKeys.all, "list", "public", username] as const,
  publicDetail: (id: string) =>
    [...ticketKeys.all, "detail", "public", id] as const,
  detail: (id: string) => [...ticketKeys.all, "detail", id] as const,
};
export function useTickets() {
  logEvent("ticket", {
    type: "get",
  });
  return useQuery({
    queryKey: ticketKeys.list(),
    queryFn: () => ticketService.getTickets(),
    staleTime: 5 * 60 * 1000,
  });
}

export function usePublicTickets(username: string) {
  logEvent("ticket", {
    type: "getpublics",
  });
  return useQuery({
    queryKey: ticketKeys.listByUsername(username),
    queryFn: () => ticketService.getPublicTickets(username),
    staleTime: 5 * 60 * 1000,
    enabled: !!username,
  });
}

export function usePublicTicket(id: string) {
  logEvent("ticket", {
    type: "getpublic",
  });
  return useQuery({
    queryKey: ticketKeys.publicDetail(id),
    queryFn: () => ticketService.getPublicTicket(id),
    staleTime: 5 * 60 * 1000,
    enabled: !!id,
  });
}

export function useTicketDetail(id: string) {
  logEvent("ticket", {
    type: "getdetail",
  });
  return useQuery({
    queryKey: ticketKeys.detail(id),
    queryFn: () => ticketService.getTicketById(id),
    enabled: !!id,
  });
}

export function useCreateTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: FormData) => ticketService.createTicket(payload),
    onSuccess: (msg) => {
      logEvent("ticket", {
        type: "create",
      });
      toast.success(msg);
      queryClient.invalidateQueries({ queryKey: ticketKeys.list() });
    },
    onError: (err) => toast.error(err.message),
  });
}

export function useUpdateTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: FormData }) =>
      ticketService.updateTicket(id, payload),
    onSuccess: (msg, variables) => {
      logEvent("ticket", {
        type: "update",
      });
      queryClient.invalidateQueries({
        queryKey: ticketKeys.detail(variables.id),
      });
      queryClient.invalidateQueries({
        queryKey: ticketKeys.publicDetail(variables.id),
      });
      toast.success(msg);
    },
    onError: (err) => toast.error(err.message),
  });
}

export function useDeleteTicket() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => ticketService.deleteTicket(id),
    onSuccess: (msg) => {
      logEvent("ticket", {
        type: "delete",
      });
      toast.success(msg);
      queryClient.invalidateQueries({ queryKey: ticketKeys.list() });
      window.location.href = "/dash/ticket";
    },
    onError: (err) => toast.error(err.message),
  });
}
