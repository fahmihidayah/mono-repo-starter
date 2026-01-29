import { useState } from "react";
import { LogOut } from "lucide-react";
import { Form } from "react-router";
import { Button } from "~/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "~/components/ui/alert-dialog";
import {
  SidebarFooter,
  SidebarMenu,
  SidebarMenuItem,
} from "~/components/ui/sidebar";

export interface UserSession {
  name: string;
  email: string;
}

interface DashboardSidebarFooterProps {
  user: UserSession;
  signOutAction?: string;
}

export function DashboardSidebarFooter({ user, signOutAction }: DashboardSidebarFooterProps) {
  const [showLogoutDialog, setShowLogoutDialog] = useState(false);

  const handleLogout = () => {
    // Submit the form programmatically
    const form = document.getElementById('logout-form') as HTMLFormElement;
    if (form) {
      form.submit();
    }
  };

  return (
    <>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <div className="flex flex-col gap-2 p-2">
              <div className="flex items-center gap-2 px-1">
                <div className="size-8 rounded-full bg-muted flex items-center justify-center">
                  <span className="text-sm font-medium">
                    {user.name?.charAt(0).toUpperCase()}
                  </span>
                </div>
                <div className="flex flex-col flex-1 min-w-0">
                  <span className="text-sm font-medium truncate">{user.name}</span>
                  <span className="text-xs text-muted-foreground truncate">
                    {user.email}
                  </span>
                </div>
              </div>

              {/* Hidden form for logout */}
              <Form
                id="logout-form"
                method="post"
                action={signOutAction || "/logout"}
              >
                <Button
                variant="outline"
                size="sm"
                className="w-full justify-start"
                // onClick={() => setShowLogoutDialog(true)}
              >
                <LogOut className="mr-2 size-4" />
                Sign out 
              </Button>
              </Form>

              {/* Logout button that shows confirmation dialog */}
              {/* <Button
                variant="outline"
                size="sm"
                className="w-full justify-start"
                onClick={() => setShowLogoutDialog(true)}
              >
                <LogOut className="mr-2 size-4" />
                Sign out aa
              </Button> */}
            </div>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>

      {/* Logout Confirmation Dialog */}
      <AlertDialog open={showLogoutDialog} onOpenChange={setShowLogoutDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Sign out of your account?</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to sign out? You'll need to sign in again to access your dashboard.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleLogout}>
              Sign out
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
