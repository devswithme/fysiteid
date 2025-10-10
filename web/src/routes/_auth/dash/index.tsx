import UserForm from "@/components/form/user";
import Modal from "@/components/modal";
import { Button, buttonVariants } from "@/components/ui/button";
import { useUpdateUser, useUser } from "@/hooks/use-user";
import type { User } from "@/types/user";
import type { userMutationType } from "@/validation/user";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/_auth/dash/")({
  component: RouteComponent,
});

function RouteComponent() {
  const [open, setOpen] = useState(false);
  const { data } = useUser();
  const { mutate, isSuccess } = useUpdateUser();

  function onSubmit(values: userMutationType) {
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
        <h1 className="text-3xl font-semibold">Profile</h1>
        <Link
          to="/$username"
          params={{ username: data?.username as string }}
          className={buttonVariants()}
        >
          View
        </Link>
      </div>
      <div className="grid grid-cols-1">
        <div className="flex flex-wrap gap-4 items-center p-4 rounded-md border justify-between">
          <div className="flex gap-4 items-center">
            <img
              src={data?.picture}
              className="size-8 rounded-full object-cover"
            />
            <div>
              <h1>{data?.name}</h1>
              <p className="text-muted-foreground">@{data?.username}</p>
            </div>
          </div>
          <Modal
            trigger={<Button variant="secondary">Edit</Button>}
            open={open}
            onOpenChange={setOpen}
          >
            <UserForm data={data as User} onSubmit={onSubmit} />
          </Modal>
        </div>
      </div>
    </>
  );
}
