import { Button } from "@/components/ui/button";
import { useCreateRegistrant } from "@/hooks/use-registrant";
import { usePublicUserByTicketId } from "@/hooks/use-user";
import { ticketService } from "@/service/ticket";
import type { Ticket } from "@/types/ticket";
import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { ImageOff } from "lucide-react";
import { useEffect } from "react";

export const Route = createFileRoute("/_public/ticket/$id/")({
  loader: async ({ params: { id } }) => {
    try {
      const data = await ticketService.getPublicTicket(id);

      if (!data) return notFound();

      return data;
    } catch {
      return notFound();
    }
  },
  validateSearch: (search: Record<string, unknown>) => {
    return {
      state: search.state as string,
      action: search.action as string,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = Route.useParams();
  const { state, action } = Route.useSearch();
  const ticket = Route.useLoaderData() as Ticket;
  const { data: user } = usePublicUserByTicketId(id);

  const { mutate } = useCreateRegistrant(state);

  useEffect(() => {
    if (action == "claim") {
      mutate(ticket.id);
    }
  }, [action, mutate, ticket]);

  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-10">
      <section className="sm:col-span-2 space-y-6">
        <div className="bg-muted aspect-video flex justify-center items-center">
          {ticket?.image ? (
            <img src={ticket.image} className="w-full h-full object-cover" />
          ) : (
            <ImageOff />
          )}
        </div>
        <div className="space-y-4">
          <h1 className="text-2xl font-semibold">{ticket?.title}</h1>
          <p>{ticket?.description}</p>
        </div>
      </section>
      <section className="border rounded-md p-4 h-fit space-y-4">
        <h1 className="font-medium text-lg">Ticket details</h1>
        <h1 className="flex justify-between">
          Quota
          <span>
            {ticket?.registered_count}/{ticket?.quota}
          </span>
        </h1>
        <Link
          to="/$username"
          params={{ username: user?.username as string }}
          className="flex gap-4 items-center border px-4 py-1 rounded-md"
        >
          <img
            src={user?.picture}
            className="size-8 rounded-full aspect-square"
          />
          <div>
            <h1 className="font-medium">{user?.name}</h1>
            <p className="text-muted-foreground text-sm">@{user?.username}</p>
          </div>
        </Link>
        <Button className="w-full" onClick={() => mutate(ticket?.id as string)}>
          Claim now
        </Button>
      </section>
    </div>
  );
}
