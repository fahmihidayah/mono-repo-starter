import { useEffect } from "react";
import { useNavigate } from "react-router";
import { Skeleton } from "~/components/ui/skeleton";
import { useSession } from "~/lib/auth-client";
import { toast } from "sonner";
import { DashboardLayout } from "~/components/layout/dashboard";

export default function DashboardRoute() {
  const navigate = useNavigate();
  const { data: session, isPending } = useSession();

  // Authentication check
  useEffect(() => {
    if (!isPending && !session) {
      toast.error("Authentication Required", {
        description: "Please sign in to access the dashboard",
      });
      navigate("/login");
    }
  }, [session, isPending, navigate]);

  // Sign out handler
  const handleSignOut = async () => {
    try {
      const authClient = (await import("~/lib/auth-client")).authClient;
      await authClient.signOut();
      toast.success("Signed out successfully");
      navigate("/");
    } catch (err) {
      console.error("Sign out failed:", err);
      toast.error("Failed to sign out");
    }
  };

  // Loading state
  if (isPending) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="space-y-4 w-64">
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-full" />
        </div>
      </div>
    );
  }

  // No session state
  if (!session) {
    return null;
  }

  // Prepare user data
  const user = {
    name: session.user.name || "User",
    email: session.user.email || "",
  };

  return (
    <DashboardLayout
      user={user}
      onSignOut={handleSignOut}
    />
  );
}
