import DataEmpty from "@/components/data-empty";
import Modal from "@/components/modal";
import { Badge } from "@/components/ui/badge";
import { useGetRegistrants } from "@/hooks/use-registrant";
import { formatDate } from "@/lib/utils";
import { createFileRoute } from "@tanstack/react-router";
import { ImageOff } from "lucide-react";
import { QRCodeSVG } from "qrcode.react";

export const Route = createFileRoute("/_auth/dash/history")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data: tickets, isLoading } = useGetRegistrants();
  return (
    <>
      <h1 className="text-3xl font-semibold">History</h1>
      <div className="grid sm:grid-cols-2 gap-6">
        {!tickets && !isLoading && <DataEmpty className="col-span-2" />}
        {tickets?.map((ticket) => (
          <Modal
            key={ticket.id}
            trigger={
              <div className="border rounded-md overflow-hidden">
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
                <div className="p-4">
                  <h1 className="break-words font-medium line-clamp-1">
                    {ticket.title}
                  </h1>
                  <Badge
                    variant={ticket.is_verified ? "success" : "destructive"}
                  >
                    {ticket.is_verified ? "Verified" : "Unverified"}
                  </Badge>

                  <Badge variant="secondary" className="mt-2">
                    {formatDate(ticket.created_at)}
                  </Badge>
                </div>
              </div>
            }
          >
            <div className="flex justify-center items-center">
              <QRCodeSVG
                value={ticket.url}
                className="w-fit h-fit p-4 bg-muted border aspect-square"
              />
              <p>{ticket.url}</p>
            </div>
          </Modal>
        ))}
      </div>
    </>
  );
}
