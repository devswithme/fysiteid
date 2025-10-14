import { buttonVariants } from "@/components/ui/button";
import { ticketKeys } from "@/hooks/use-ticket";
import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Check, X } from "lucide-react";

export const Route = createFileRoute("/_public/verify/")({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      err: search.err as string,
      ticket: search.ticket as string,
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { err, ticket } = Route.useSearch();
  const queryClient = useQueryClient();

  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <div className="flex flex-col items-center space-y-4 text-center">
        {err == "1" ? (
          <>
            <X className="w-16 h-16 text-destructive" />
            <h1 className="text-2xl font-semibold">Failed to verify ticket</h1>
            <Link to="/" className={buttonVariants()}>
              Home
            </Link>
          </>
        ) : (
          <>
            <Check className="w-16 h-16 text-primary" />
            <h1 className="text-2xl font-semibold">
              Successfully verified ticket
            </h1>
            {/* @ts-expect-error search is required */}
            <Link
              to="/dash/ticket/$id"
              params={{ id: ticket }}
              onClick={() => {
                queryClient.invalidateQueries({
                  queryKey: ticketKeys.publicDetail(ticket),
                });
              }}
              className={buttonVariants()}
            >
              Go back
            </Link>
          </>
        )}
      </div>
    </div>
  );
}
