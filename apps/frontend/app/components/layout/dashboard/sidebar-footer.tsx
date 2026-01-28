import { LogOut } from "lucide-react";
import { Button } from "~/components/ui/button";
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
  onSignOut: () => void;
}

export function DashboardSidebarFooter({ user, onSignOut }: DashboardSidebarFooterProps) {
  return (
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
            <Button
              variant="outline"
              size="sm"
              onClick={onSignOut}
              className="w-full justify-start"
            >
              <LogOut className="mr-2 size-4" />
              Sign out
            </Button>
          </div>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>
  );
}
