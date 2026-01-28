import { Home, LayoutDashboard, Users, Settings, Package, CheckSquare } from "lucide-react";
import type { NavigationGroup } from "./sidebar-navigation";
import type { SidebarHeaderConfig } from "./sidebar-header";

// Default header configuration
export const defaultHeaderConfig: SidebarHeaderConfig = {
  appName: "Starter App",
  appInitial: "S",
  subtitle: "Dashboard",
};

// Default navigation groups
export const defaultNavigationGroups: NavigationGroup[] = [
  {
    label: "Navigation",
    items: [
      {
        title: "Overview",
        url: "/dashboard",
        icon: LayoutDashboard,
      },
      {
        title: "Tasks",
        url: "/dashboard/tasks",
        icon: CheckSquare,
      },
      {
        title: "Users",
        url: "/dashboard/users",
        icon: Users,
      },
      {
        title: "Settings",
        url: "/dashboard/settings",
        icon: Settings,
      },
    ],
  },
  {
    label: "Quick Links",
    items: [
      {
        title: "Home",
        url: "/",
        icon: Home,
      },
    ],
  },
];
