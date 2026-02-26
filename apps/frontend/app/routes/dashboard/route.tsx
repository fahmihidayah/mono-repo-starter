import { useLoaderData } from "react-router";
import { useEffect } from "react";
import { toast } from "sonner";
import { requireAuth } from "~/lib/auth-guard.server";
import { DashboardLayout } from "~/components/layout/dashboard";
import type { Route } from "./+types/route";
import { FileText, FolderTree, Home, List, Settings, User, Users } from "lucide-react";
import * as serverSession from "~/session.server";

/**
 * Loader - Server-side authentication check
 * Runs before component renders, redirects to /login if not authenticated
 */
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);
  const { getSession } = serverSession;
  // Get flash messages from session
  const session = await getSession(request.headers.get("Cookie"));
  const error = session.get("error");
  const success = session.get("success");

  return {
    user,
    flash: { error, success },
  } as const;
}

/**
 * Dashboard Route Component
 * User is guaranteed to be authenticated when this renders
 */
export default function DashboardRoute() {
  const { user, flash } = useLoaderData<typeof loader>();

  // Show flash messages as toasts
  useEffect(() => {
    if (flash?.error) {
      toast.error(flash.error);
    }
    if (flash?.success) {
      toast.success(flash.success);
    }
  }, [flash]);

  return (
    <DashboardLayout
      user={{
        name: user.userName,
        email: user.userEmail,
      }}
      config={{
        header: {
          appInitial: "SA",
          appName: "Starter Apps",
          subtitle: "Starter app",
        },
        headerTitle: "Main",
        navigationGroups: [
          {
            label: "Main",
            items: [
              {
                title: "Home",
                url: "/dashboard",
                icon: Home,
              },
              {
                title: "Media",
                url: "/dashboard/media",
                icon: FileText,
              },
              {
                title: "Users",
                url: "/dashboard/users",
                icon: Users,
              },
              {
                title: "Categories",
                url: "/dashboard/categories",
                icon: FolderTree,
              },
              {
                title: "Post",
                url: "/dashboard/posts",
                icon: FileText,
              },
            ],
          },
          {
            label: "Setting",
            items: [
              {
                title: "Settings",
                url: "/dashboard/setting",
                icon: Settings,
              },
            ],
          },
        ],
      }}
      signOutAction="/logout"
    />
  );
}
