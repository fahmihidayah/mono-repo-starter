import { Link, useLocation } from "react-router";
import {
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "~/components/ui/sidebar";
import { Separator } from "~/components/ui/separator";
import type { LucideIcon } from "lucide-react";

export interface NavigationItem {
  title: string;
  url: string;
  icon: LucideIcon;
}

export interface NavigationGroup {
  label: string;
  items: NavigationItem[];
}

interface DashboardSidebarNavigationProps {
  groups: NavigationGroup[];
}

export function DashboardSidebarNavigation({ groups }: DashboardSidebarNavigationProps) {
  const location = useLocation();

  return (
    <SidebarContent>
      {groups.map((group, index) => (
        <div key={group.label}>
          {index > 0 && <Separator className="my-2" />}
          <SidebarGroup>
            <SidebarGroupLabel>{group.label}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      asChild
                      isActive={location.pathname === item.url}
                      tooltip={item.title}
                    >
                      <Link to={item.url}>
                        <item.icon />
                        <span>{item.title}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </div>
      ))}
    </SidebarContent>
  );
}
