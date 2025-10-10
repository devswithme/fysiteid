import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useLogout } from "@/hooks/use-auth";
import { useUser } from "@/hooks/use-user";
import { items } from "@/lib/const/routes";
import { LogOut } from "lucide-react";

export function AppSidebar() {
  const { data: user, isLoading } = useUser();
  const { mutate: logout } = useLogout();

  return (
    <Sidebar>
      <SidebarContent className="bg-background">
        <SidebarGroup className="!p-0">
          <SidebarGroupContent>
            <SidebarMenu className="gap-0">
              {!isLoading && (
                <SidebarMenuItem className="p-4 flex gap-2 items-center">
                  <img
                    src={user?.picture}
                    referrerPolicy="no-referrer"
                    className="size-6 rounded-full object-cover"
                  />
                  <h1>{user?.name.split(" ")[0]}</h1>
                </SidebarMenuItem>
              )}
              <div className="px-2">
                {items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild>
                      <a href={item.url}>
                        <item.icon />
                        <span>{item.title}</span>
                      </a>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
                <SidebarMenuItem>
                  <SidebarMenuButton
                    className="cursor-pointer"
                    onClick={() => logout()}
                  >
                    <LogOut /> Logout
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </div>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  );
}
