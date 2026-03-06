import { Link, useLocation } from "react-router";
import { useState } from "react";
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
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "~/components/ui/collapsible";
import { ChevronDown } from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { cn } from "~/lib/utils";

export interface NavigationItem {
  title: string;
  url: string;
  icon: LucideIcon;
}

export interface NavigationGroup {
  label: string;
  items: NavigationItem[];
  collapsible?: boolean;
  defaultOpen?: boolean;
}

interface DashboardSidebarNavigationProps {
  groups: NavigationGroup[];
}

export function DashboardSidebarNavigation({ groups }: DashboardSidebarNavigationProps) {
  const location = useLocation();
  const [openGroups, setOpenGroups] = useState<Record<string, boolean>>(() => {
    const initial: Record<string, boolean> = {};
    groups.forEach((group) => {
      initial[group.label] = group.defaultOpen ?? true;
    });
    return initial;
  });

  const toggleGroup = (label: string) => {
    setOpenGroups((prev) => ({ ...prev, [label]: !prev[label] }));
  };

  return (
    <SidebarContent>
      {groups.map((group, index) => {
        const isCollapsible = group.collapsible ?? false;
        const isOpen = openGroups[group.label] ?? true;

        return (
          <div key={group.label}>
            {index > 0 && <Separator className="my-2" />}
            <SidebarGroup>
              {isCollapsible ? (
                <Collapsible open={isOpen} onOpenChange={() => toggleGroup(group.label)}>
                  <CollapsibleTrigger className="w-full">
                    <SidebarGroupLabel className="flex items-center justify-between cursor-pointer group">
                      <span>{group.label}</span>
                      <ChevronDown
                        className={cn(
                          "h-4 w-4 transition-transform duration-200",
                          isOpen && "rotate-180"
                        )}
                      />
                    </SidebarGroupLabel>
                  </CollapsibleTrigger>
                  <CollapsibleContent>
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
                  </CollapsibleContent>
                </Collapsible>
              ) : (
                <>
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
                </>
              )}
            </SidebarGroup>
          </div>
        );
      })}
    </SidebarContent>
  );
}
