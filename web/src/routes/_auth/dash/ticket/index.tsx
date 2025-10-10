import Modal from "@/components/modal";

import { useCreateTicket, useTickets } from "@/hooks/use-ticket";
import { createFileRoute } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import TicketForm from "@/components/form/ticket";
import type { ticketMutationType } from "@/validation/ticket";
import DataEmpty from "@/components/data-empty";
import { useEffect, useState } from "react";
import { Card } from "@/components/card";

export const Route = createFileRoute("/_auth/dash/ticket/")({
  component: RouteComponent,
});

function RouteComponent() {
  const [open, setOpen] = useState(false);
  const { data, isLoading } = useTickets();
  const { mutate, isSuccess } = useCreateTicket();

  function onSubmit(values: ticketMutationType) {
    const formData = new FormData();

    Object.entries(values).forEach(([key, value]) => {
      formData.append(key, value as string);
    });

    mutate(formData);
  }

  useEffect(() => {
    if (isSuccess) setOpen(false);
  }, [isSuccess]);

  return (
    <>
      <div className="flex justify-between flex-wrap gap-4">
        <h1 className="text-3xl font-semibold">Ticket</h1>
        <Modal
          trigger={<Button>Add</Button>}
          open={open}
          onOpenChange={setOpen}
        >
          <TicketForm onSubmit={onSubmit} />
        </Modal>
      </div>
      <div className="grid sm:grid-cols-2 gap-6">
        {!data && !isLoading && <DataEmpty className="col-span-2" />}
        {data?.map((ticket) => (
          <Card ticket={ticket} key={ticket.id} />
        ))}
      </div>
    </>
  );
}
