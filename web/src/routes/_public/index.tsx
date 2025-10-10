import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_public/")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <>
      <h1 className="text-3xl sm:text-6xl font-semibold sm:text-center tracking-tight leading-tighter text-balance">
        Simple yet reliable ticketing system
      </h1>
      <div className="bg-neutral-100 sm:p-8 rounded-3xl">
        <img src="/home.png" />
      </div>
    </>
  );
}
