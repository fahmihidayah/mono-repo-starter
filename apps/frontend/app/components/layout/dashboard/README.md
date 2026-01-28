# Dashboard Layout Components

A modular and configurable dashboard layout system with separate, reusable components.

## Structure

```
dashboard/
├── index.tsx                    # Main exports
├── dashboard-layout.tsx         # Main layout wrapper
├── sidebar-header.tsx           # Sidebar header component
├── sidebar-navigation.tsx       # Sidebar navigation component
├── sidebar-footer.tsx           # Sidebar footer with user info
├── config.ts                    # Default configuration
└── README.md                    # This file
```

## Components

### DashboardLayout

The main layout component that combines all sidebar components.

**Props:**
- `user: UserSession` - User information (name, email)
- `onSignOut: () => void` - Sign out handler
- `config?: DashboardLayoutConfig` - Optional configuration

### DashboardSidebarHeader

Displays the app name and logo in the sidebar header.

### DashboardSidebarNavigation

Renders navigation groups with links and icons.

### DashboardSidebarFooter

Shows user profile and sign out button.

## Usage

### Basic Usage (Default Configuration)

```tsx
import { DashboardLayout } from "~/components/layout/dashboard";

export default function DashboardRoute() {
  const user = {
    name: "John Doe",
    email: "john@example.com",
  };

  const handleSignOut = async () => {
    // Sign out logic
  };

  return (
    <DashboardLayout
      user={user}
      onSignOut={handleSignOut}
    />
  );
}
```

### Custom Configuration

```tsx
import { DashboardLayout } from "~/components/layout/dashboard";
import { Building, FileText, Settings } from "lucide-react";

export default function DashboardRoute() {
  const user = {
    name: "John Doe",
    email: "john@example.com",
  };

  const customConfig = {
    header: {
      appName: "My Business App",
      appInitial: "B",
      subtitle: "Admin Panel",
    },
    navigationGroups: [
      {
        label: "Main",
        items: [
          {
            title: "Dashboard",
            url: "/dashboard",
            icon: Building,
          },
          {
            title: "Documents",
            url: "/dashboard/documents",
            icon: FileText,
          },
        ],
      },
      {
        label: "System",
        items: [
          {
            title: "Settings",
            url: "/dashboard/settings",
            icon: Settings,
          },
        ],
      },
    ],
    headerTitle: "Admin Dashboard",
  };

  return (
    <DashboardLayout
      user={user}
      onSignOut={handleSignOut}
      config={customConfig}
    />
  );
}
```

## Configuration Options

### SidebarHeaderConfig

```typescript
interface SidebarHeaderConfig {
  appName: string;      // App name to display
  appInitial: string;   // Single letter for logo
  subtitle?: string;    // Optional subtitle (default: "Dashboard")
}
```

### NavigationGroup

```typescript
interface NavigationGroup {
  label: string;        // Group label
  items: NavigationItem[];
}

interface NavigationItem {
  title: string;        // Link text
  url: string;          // Route URL
  icon: LucideIcon;     // Lucide icon component
}
```

### DashboardLayoutConfig

```typescript
interface DashboardLayoutConfig {
  header?: SidebarHeaderConfig;       // Sidebar header config
  navigationGroups?: NavigationGroup[]; // Navigation groups
  headerTitle?: string;                // Main header title
}
```

## Why Props Instead of Context API?

This implementation uses **props** instead of Context API for several reasons:

### Advantages of Props:

1. **Simplicity**: Props are straightforward and easy to understand
2. **Explicit Data Flow**: Clear parent-to-child relationship
3. **Better Performance**: No context re-renders
4. **Type Safety**: Full TypeScript support without extra setup
5. **Testing**: Easier to test components in isolation
6. **Flexibility**: Each route can have different configs
7. **No Provider Wrapper**: Less boilerplate code

### When to Use Context API Instead:

Consider using Context API if you need:
- Deep component nesting (3+ levels)
- Global state shared across many routes
- Frequent config changes at runtime
- Multiple consumers at different nesting levels

For this dashboard layout, props provide the best balance of simplicity and functionality.

## Customizing Default Configuration

Edit [config.ts](config.ts) to change the default settings:

```typescript
// app/components/layout/dashboard/config.ts
export const defaultHeaderConfig: SidebarHeaderConfig = {
  appName: "Your App Name",
  appInitial: "Y",
  subtitle: "Dashboard",
};

export const defaultNavigationGroups: NavigationGroup[] = [
  // Your navigation groups
];
```

## Examples

### Different Configs for Different Routes

```tsx
// Admin dashboard
<DashboardLayout
  config={{
    header: { appName: "Admin Panel", appInitial: "A" },
    headerTitle: "Admin Dashboard",
  }}
/>

// User dashboard
<DashboardLayout
  config={{
    header: { appName: "My App", appInitial: "M" },
    headerTitle: "My Dashboard",
  }}
/>
```

### Adding Custom Icons

```tsx
import { Rocket, Zap, Star } from "lucide-react";

const customNav = [
  {
    label: "Features",
    items: [
      { title: "Launch", url: "/launch", icon: Rocket },
      { title: "Performance", url: "/perf", icon: Zap },
      { title: "Favorites", url: "/favorites", icon: Star },
    ],
  },
];
```
