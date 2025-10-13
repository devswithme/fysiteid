import { Link } from "@tanstack/react-router";
import { Badge } from "./ui/badge";
import type { Ticket } from "@/types/ticket";
import { ImageOff } from "lucide-react";
import { Route as TicketRoute } from "@/routes/_auth/dash/ticket/$id/index";
export function Card({ ticket }: { ticket: Ticket }) {
  return (
    <Link
      to={TicketRoute.to}
      params={{ id: ticket.id }}
      search={{ page: 1 }}
      className="border rounded-md overflow-hidden"
    >
      <div className="bg-muted aspect-video flex justify-center items-center">
        {ticket.image ? (
          <img src={ticket.image} className="w-full h-full object-cover" />
        ) : (
          <ImageOff />
        )}
      </div>
      <div className="p-4">
        <h1 className="break-words font-medium line-clamp-1">{ticket.title}</h1>
        <p className="break-words line-clamp-1 text-sm text-muted-foreground">
          {ticket.description}
        </p>
        <Badge variant="secondary" className="mt-2">
          {ticket.registered_count}/{ticket.quota}
        </Badge>
      </div>
    </Link>
  );
}
