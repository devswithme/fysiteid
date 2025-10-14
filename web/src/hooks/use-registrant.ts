import { logEvent } from "@/lib/utils";
import { authService } from "@/service/auth";
import { registrantService } from "@/service/registrant";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";

export const registrantKeys = {
  all: ["registrant"] as const,
  list: () => [...registrantKeys.all, "list"] as const,
  detail: (id: string) => [...registrantKeys.all, "detail", id] as const,
  gen: (id: string) => [...registrantKeys.all, "gen", id] as const,
};

export function useCreateRegistrant(state: string) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (ticketID: string) =>
      registrantService.createRegistrant(ticketID, state),
    onError: (err, ticketID) => {
      if (err.message === "Unauthorized") {
        authService.login(`/ticket/${ticketID}?state=${state}&action=claim`);
      } else {
        toast.error("Fail to claim ticket");
      }
    },
    onSuccess: (_, ticketID) => {
      logEvent("registrant", {
        type: "create",
      });
      queryClient.invalidateQueries({
        queryKey: registrantKeys.detail(ticketID),
      });
      queryClient.invalidateQueries({ queryKey: registrantKeys.list() });
      toast.success("Successfully claim ticket");
      navigate({ to: "/dash/history" });
    },
  });
}

export function useGenerateState(ticketID: string) {
  logEvent("registrant", {
    type: "state",
  });
  return useQuery({
    queryKey: registrantKeys.gen(ticketID),
    queryFn: () => registrantService.generateState(ticketID),
    enabled: !!ticketID,
    refetchOnWindowFocus: false,
  });
}

export function useGetRegistrantByTicketID({
  ticketID,
  page,
  limit,
  search,
}: {
  ticketID: string;
  page?: number;
  limit?: number;
  search: string;
}) {
  logEvent("registrant", {
    type: "getbyticket",
  });
  return useQuery({
    queryKey: [...registrantKeys.detail(ticketID), { page, limit, search }],
    queryFn: () =>
      registrantService.getRegistrantByTicketID(ticketID, page, limit, search),
    staleTime: 5 * 60 * 1000,
    enabled: !!ticketID,
  });
}

export function useGetRegistrants() {
  logEvent("registrant", {
    type: "getbyuser",
  });
  return useQuery({
    queryKey: registrantKeys.list(),
    queryFn: () => registrantService.getRegistrantByUserID(),
    staleTime: 5 * 60 * 1000,
  });
}
