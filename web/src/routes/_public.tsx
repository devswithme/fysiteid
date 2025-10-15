import { Button, buttonVariants } from "@/components/ui/button";
import { useLogin } from "@/hooks/use-auth";
import { useUser } from "@/hooks/use-user";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_public")({
  component: RouteComponent,
});

function RouteComponent() {
  const { login } = useLogin();
  const { data: user } = useUser();
  return (
    <main>
      <header className="w-full">
        <nav className="max-w-3xl mx-auto flex justify-between items-center px-4 sm:px-0 py-6">
          <Link to="/">
            <img src="/logo.svg" className="w-20" />
          </Link>
          {user ? (
            <Link to="/dash" className={buttonVariants()}>
              Dash
            </Link>
          ) : (
            <Button onClick={login}>Login</Button>
          )}
        </nav>
      </header>
      <main className="flex flex-col justify-center items-center gap-y-16 py-12 max-w-3xl mx-auto px-4 sm:px-0">
        <Outlet />
      </main>
    </main>
  );
}
