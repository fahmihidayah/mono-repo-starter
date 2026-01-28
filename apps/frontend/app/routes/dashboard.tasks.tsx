import { Outlet } from "react-router";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "~/components/ui/card";
import { useSession } from "~/lib/auth-client";

export function meta() {
  return [
    { title: "Dashboard - Starter App" },
    { name: "description", content: "Your personal dashboard" },
  ];
}

export default function DashboardIndex() {
  const { data: session } = useSession();

  if (!session) {
    return null;
  }

  return (
    <Outlet />
  );
}
