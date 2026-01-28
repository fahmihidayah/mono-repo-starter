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
    <div className="flex-1 p-6">
      <div className="space-y-6">
        {/* Welcome Section */}
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            Welcome back, {session.user.name}!
          </h1>
          <p className="text-muted-foreground mt-2">
            Here's what's happening with your account today.
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid gap-6 md:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Account Status</CardTitle>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                className="h-4 w-4 text-muted-foreground"
              >
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                <circle cx="9" cy="7" r="4" />
                <path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
              </svg>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {session.user.emailVerified ? "Verified" : "Unverified"}
              </div>
              <p className="text-xs text-muted-foreground">
                {session.user.emailVerified
                  ? "Your email is verified"
                  : "Please verify your email"}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Email</CardTitle>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                className="h-4 w-4 text-muted-foreground"
              >
                <rect width="20" height="16" x="2" y="4" rx="2" />
                <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
              </svg>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold truncate">{session.user.email}</div>
              <p className="text-xs text-muted-foreground">Your account email</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Member Since</CardTitle>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                className="h-4 w-4 text-muted-foreground"
              >
                <path d="M8 2v4" />
                <path d="M16 2v4" />
                <rect width="18" height="18" x="3" y="4" rx="2" />
                <path d="M3 10h18" />
              </svg>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {new Date(session.user.createdAt).toLocaleDateString()}
              </div>
              <p className="text-xs text-muted-foreground">Account creation date</p>
            </CardContent>
          </Card>
        </div>

        {/* User Details Card */}
        <Card>
          <CardHeader>
            <CardTitle>Account Details</CardTitle>
            <CardDescription>
              Your account information and settings
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <p className="text-sm font-medium text-muted-foreground">User ID</p>
                <p className="text-sm font-mono mt-1">{session.user.id}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Full Name</p>
                <p className="text-sm mt-1">{session.user.name}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Email Address</p>
                <p className="text-sm mt-1">{session.user.email}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">
                  Email Verification
                </p>
                <p className="text-sm mt-1">
                  {session.user.emailVerified ? (
                    <span className="text-green-600 dark:text-green-400">✓ Verified</span>
                  ) : (
                    <span className="text-yellow-600 dark:text-yellow-400">
                      ⚠ Not verified
                    </span>
                  )}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Activity Card - Example of easy modification */}
        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>Your latest account activity</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center justify-between py-2">
                <div className="flex items-center gap-3">
                  <div className="size-8 rounded-full bg-primary/10 flex items-center justify-center">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      className="size-4 text-primary"
                    >
                      <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
                      <polyline points="10 17 15 12 10 7" />
                      <line x1="15" x2="3" y1="12" y2="12" />
                    </svg>
                  </div>
                  <div>
                    <p className="text-sm font-medium">Signed in</p>
                    <p className="text-xs text-muted-foreground">Successfully logged in</p>
                  </div>
                </div>
                <span className="text-xs text-muted-foreground">Just now</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
