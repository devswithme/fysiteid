import DataEmpty from "@/components/data-empty";
import { usePublicTickets } from "@/hooks/use-ticket";
import { userService } from "@/service/user";
import type { User } from "@/types/user";
import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { ImageOff } from "lucide-react";

export const Route = createFileRoute("/_public/$username")({
  loader: async ({ params: { username } }) => {
    try {
      const data = await userService.getPublicUser(username);

      if (!data) return notFound();

      return data;
    } catch {
      return notFound();
    }
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { username } = Route.useParams();
  const user = Route.useLoaderData() as User;
  const { data: tickets } = usePublicTickets(username);
  return (
    <>
      <div className="flex flex-col justify-center items-center gap-y-4 text-center">
        <img
          src={user?.picture}
          className="size-18 rounded-full object-cover"
        />
        <div>
          <h1 className="font-medium">{user?.name}</h1>
          <p className="text-muted-foreground">@{user?.username}</p>
        </div>
      </div>
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 w-full">
        {!tickets?.length && <DataEmpty className="col-span-2" />}
        {tickets?.map((ticket) => (
          // @ts-expect-error property search is required
          <Link
            to="/ticket/$id"
            params={{ id: ticket.id }}
            key={ticket.id}
            className="space-y-3"
          >
            <div className="bg-muted aspect-video flex justify-center items-center">
              {ticket.image ? (
                <img
                  src={ticket.image}
                  className="w-full h-full object-cover"
                />
              ) : (
                <ImageOff />
              )}
            </div>
            <h1 className="font-medium line-clamp-2">{ticket.title}</h1>
          </Link>
        ))}
      </div>
    </>
  );
}
