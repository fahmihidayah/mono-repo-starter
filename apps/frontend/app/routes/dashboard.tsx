import { useLoaderData } from "react-router";
import { requireAuth } from "~/lib/auth-guard.server";
import { DashboardLayout } from "~/components/layout/dashboard";
import type { Route } from "./+types/dashboard";

/**
 * Loader - Server-side authentication check
 * Runs before component renders, redirects to /login if not authenticated
 */
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);
  return { user };
}

/**
 * Dashboard Route Component
 * User is guaranteed to be authenticated when this renders
 */
export default function DashboardRoute() {
  const { user } = useLoaderData<typeof loader>();

  return (
    <DashboardLayout
      user={{
        name: user.userName,
        email: user.userEmail,
      }}
      onSignOut={() => {
        // Use Form to submit to logout action
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/logout';
        document.body.appendChild(form);
        form.submit();
      }}
    />
  );
}
