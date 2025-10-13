"use client";

import { Edit, MoreHorizontal, Share, Trash } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Sidebar,
  SidebarContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useEffect, useState } from "react";
import { useDeleteTicket, useUpdateTicket } from "@/hooks/use-ticket";
import type { ticketMutationType } from "@/validation/ticket";
import Modal from "../modal";
import TicketForm from "../form/ticket";
import { useMatch } from "@tanstack/react-router";

export function NavActions() {
  const match = useMatch({
    from: "/_auth/dash/ticket/$id/",
    shouldThrow: false,
  });

  const id = match?.params?.id as string;
  const ticket = match ? match.loaderData : null;

  const [isOpen, setIsOpen] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const { mutate: bin } = useDeleteTicket();
  const { mutate: update, isSuccess } = useUpdateTicket();

  async function onSubmit(values: ticketMutationType) {
    const formData = new FormData();

    Object.entries(values).forEach(([key, value]) => {
      formData.append(key, value as string);
    });

    update({ id, payload: formData });
  }

  useEffect(() => {
    if (isSuccess) setIsModalOpen(false);
  }, [isSuccess]);

  return (
    <div className="flex items-center gap-2 text-sm">
      <Button
        variant="ghost"
        size="icon"
        className="h-7 w-7"
        onClick={() => {
          if (confirm("Are you sure? this action is not revertable")) bin(id);
        }}
      >
        <Trash />
      </Button>
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className="data-[state=open]:bg-accent h-7 w-7"
          >
            <MoreHorizontal />
          </Button>
        </PopoverTrigger>
        <PopoverContent
          className="w-56 overflow-hidden rounded-lg p-0"
          align="end"
        >
          <Sidebar collapsible="none" className="bg-transparent">
            <SidebarContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <Modal
                    trigger={
                      <SidebarMenuButton>
                        <Edit />
                        Update
                      </SidebarMenuButton>
                    }
                    open={isModalOpen}
                    onOpenChange={setIsModalOpen}
                  >
                    <TicketForm onSubmit={onSubmit} data={ticket!} />
                  </Modal>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton
                    variant="outline"
                    onClick={async () => {
                      await navigator.clipboard.writeText(
                        `${import.meta.env.VITE_APP_URL}/ticket/${id}`
                      );
                    }}
                  >
                    <Share />
                    Share
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarContent>
          </Sidebar>
        </PopoverContent>
      </Popover>
    </div>
  );
}
