import { AppSidebar } from "@/components/sidebar/app";
import Navigation from "@/components/nav";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { authMiddleware } from "@/lib/middleware";
import { createFileRoute, Outlet, useLocation } from "@tanstack/react-router";
import { Separator } from "@/components/ui/separator";
import { NavActions } from "@/components/sidebar/action";
import { items } from "@/lib/const/routes";

export const Route = createFileRoute("/_auth")({
  beforeLoad: async () => await authMiddleware.requireAuth(),
  component: RouteComponent,
});

function RouteComponent() {
  const { pathname } = useLocation();

  return (
    <>
      <SidebarProvider>
        <AppSidebar />
        <main className="w-full">
          <SidebarInset>
            <header className="flex h-14 shrink-0 items-center gap-2">
              <div className="flex flex-1 items-center gap-2 px-3">
                <SidebarTrigger />
                <Separator
                  orientation="vertical"
                  className="mr-2 data-[orientation=vertical]:h-4"
                />
                <Navigation />
              </div>
              {!items.some((item) => item.url === pathname) && (
                <div className="ml-auto px-3">
                  <NavActions />
                </div>
              )}
            </header>
          </SidebarInset>
          <section className="max-w-3xl mx-auto py-16 px-6 space-y-16">
            <Outlet />
          </section>
        </main>
      </SidebarProvider>
    </>
  );
}
