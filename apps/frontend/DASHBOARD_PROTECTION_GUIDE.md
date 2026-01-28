# Dashboard Protection Implementation

## Overview
The dashboard routes now use proper React Router server-side authentication with loaders, following the official React Router patterns for route protection.

---

## What Changed

### Before (Client-Side Auth Check)

**Problems**:
1. ❌ Client-side auth check with `useSession()` hook
2. ❌ Flash of unauthenticated content
3. ❌ Security risk (client-side checks can be bypassed)
4. ❌ Multiple useEffect dependencies and loading states
5. ❌ Manual navigation logic

**Code**:
```typescript
export default function DashboardRoute() {
  const { data: session, isPending } = useSession();

  useEffect(() => {
    if (!isPending && !session) {
      navigate("/login");
    }
  }, [session, isPending, navigate]);

  if (isPending) return <LoadingSpinner />;
  if (!session) return null;

  return <Dashboard user={session.user} />;
}
```

### After (Server-Side Auth with Loader)

**Solutions**:
1. ✅ Server-side auth check in loader
2. ✅ No flash of unauthenticated content
3. ✅ Secure (runs on server before rendering)
4. ✅ Simple component logic
5. ✅ Automatic redirect with React Router

**Code**:
```typescript
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);
  return { user };
}

export default function DashboardRoute() {
  const { user } = useLoaderData<typeof loader>();
  return <Dashboard user={user} />;
}
```

---

## Architecture

### File Structure

```
apps/frontend/app/
├── session.server.ts                    # Session storage
├── lib/
│   ├── auth.server.ts                   # Auth API calls
│   ├── auth-guard.server.ts             # Route protection (requireAuth)
│   └── auth-client.ts                   # Client utilities (signOut)
└── routes/
    ├── logout.tsx                       # Logout action (NEW)
    ├── dashboard.tsx                    # Protected layout (ENHANCED)
    └── dashboard._index.tsx             # Protected page (ENHANCED)
```

---

## Key Components

### 1. Logout Route (`routes/logout.tsx`)

**Purpose**: Handle user logout with proper session destruction

**Features**:
- Destroys server-side session
- Calls API logout endpoint
- Redirects to login page
- POST-only for security (GET requests redirect)

**Code**:
```typescript
export async function action({ request }: Route.ActionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const token = session.get("token");

  // Call API logout
  if (token) {
    await callLogoutApi(token);
  }

  // Destroy session and redirect
  return redirect("/login", {
    headers: { "Set-Cookie": await destroySession(session) }
  });
}

export async function loader() {
  return redirect("/"); // GET requests not allowed
}
```

**Usage**:
```typescript
// From anywhere in the app
<Form method="post" action="/logout">
  <button type="submit">Sign Out</button>
</Form>

// Or programmatically
const form = document.createElement('form');
form.method = 'POST';
form.action = '/logout';
document.body.appendChild(form);
form.submit();
```

---

### 2. Protected Dashboard Layout (`routes/dashboard.tsx`)

**Purpose**: Main dashboard layout with authentication

**Features**:
- Server-side auth check in loader
- User data from session
- Sign out handler

**Code**:
```typescript
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);
  return { user };
}

export default function DashboardRoute() {
  const { user } = useLoaderData<typeof loader>();

  return (
    <DashboardLayout
      user={{
        name: user.userName,
        email: user.userEmail,
      }}
      onSignOut={() => {
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/logout';
        document.body.appendChild(form);
        form.submit();
      }}
    />
  );
}
```

---

### 3. Protected Dashboard Pages (`routes/dashboard._index.tsx`)

**Purpose**: Dashboard homepage with user data

**Features**:
- Inherits parent loader protection
- Can add own loader for additional data
- Uses loader data for rendering

**Code**:
```typescript
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);
  return { user };
}

export default function DashboardIndex() {
  const { user } = useLoaderData<typeof loader>();

  return (
    <div>
      <h1>Welcome back, {user.userName}!</h1>
      <p>Email: {user.userEmail}</p>
      <p>User ID: {user.userId}</p>
    </div>
  );
}
```

---

## Authentication Flow

### Login Flow
```
1. User submits login form
   ↓
2. Action validates credentials
   ↓
3. API returns auth data (token, exp, user info)
   ↓
4. Create session with user data
   ↓
5. Calculate maxAge from token exp
   ↓
6. Commit session cookie
   ↓
7. Redirect to /dashboard
   ↓
8. Dashboard loader runs
   ↓
9. requireAuth() checks session
   ↓
10. Session exists → Return user data
   ↓
11. Dashboard renders with user data
```

### Protected Route Access
```
User navigates to /dashboard
   ↓
Dashboard loader runs (server-side)
   ↓
requireAuth(request) called
   ↓
Check session cookie
   ↓
Session exists?
   ├─ Yes → Return user data → Render dashboard
   └─ No → Throw redirect("/login")
```

### Logout Flow
```
User clicks "Sign Out"
   ↓
POST request to /logout
   ↓
Logout action runs
   ↓
Get session from cookie
   ↓
Call API logout endpoint
   ↓
Destroy session
   ↓
Redirect to /login with destroyed session cookie
   ↓
User logged out
```

---

## Benefits

### 1. **Security**
- ✅ Server-side authentication (can't be bypassed)
- ✅ No client-side auth logic to circumvent
- ✅ Session checked before component renders
- ✅ Automatic redirect on unauthorized access

### 2. **User Experience**
- ✅ No flash of unauthenticated content
- ✅ No loading spinners for auth checks
- ✅ Instant redirect if not authenticated
- ✅ Smooth navigation

### 3. **Developer Experience**
- ✅ Simple component logic
- ✅ No useEffect dependencies
- ✅ Type-safe loader data
- ✅ Easy to add protected routes

### 4. **Performance**
- ✅ Auth check happens during navigation
- ✅ No extra render cycles
- ✅ No client-side redirect delay
- ✅ Efficient server-side checks

---

## Adding New Protected Routes

### Method 1: Add Loader to Existing Route

```typescript
// app/routes/dashboard.settings.tsx
import { useLoaderData } from "react-router";
import { requireAuth } from "~/lib/auth-guard.server";
import type { Route } from "./+types/dashboard.settings";

export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);
  return { user };
}

export default function Settings() {
  const { user } = useLoaderData<typeof loader>();
  return <div>Settings for {user.userName}</div>;
}
```

### Method 2: Inherit from Parent Loader

If your route is under `/dashboard/*`, it can inherit auth from parent:

```typescript
// app/routes/dashboard.profile.tsx
import { useLoaderData } from "react-router";
import type { Route } from "./+types/dashboard.profile";

// No loader needed - inherits from dashboard.tsx
export default function Profile() {
  // Access parent loader data
  return <div>Profile Page</div>;
}
```

### Method 3: Optional Authentication

For routes that work with or without auth:

```typescript
import { getAuthUser } from "~/lib/auth-guard.server";

export async function loader({ request }: Route.LoaderArgs) {
  const user = await getAuthUser(request); // Returns null if not authenticated
  return { user };
}

export default function OptionalAuthRoute() {
  const { user } = useLoaderData<typeof loader>();

  return (
    <div>
      {user ? (
        <p>Welcome back, {user.userName}!</p>
      ) : (
        <p>Please <Link to="/login">sign in</Link></p>
      )}
    </div>
  );
}
```

---

## User Data Available in Loaders

From `requireAuth()` or `getAuthUser()`:

```typescript
interface AuthenticatedUser {
  userId: string;      // User's unique ID
  token: string;       // JWT token for API calls
  userName: string;    // User's display name
  userEmail: string;   // User's email address
}
```

**Usage**:
```typescript
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);

  // Make authenticated API calls
  const response = await fetch(`${API_URL}/api/users/${user.userId}/data`, {
    headers: {
      'Authorization': `Bearer ${user.token}`,
    },
  });

  return { user, data: await response.json() };
}
```

---

## Common Patterns

### 1. Fetch User-Specific Data

```typescript
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);

  // Fetch data for this specific user
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
  const response = await fetch(`${apiBaseUrl}/api/users/${user.userId}/profile`, {
    headers: {
      'Authorization': `Bearer ${user.token}`,
    },
  });

  const profile = await response.json();
  return { user, profile };
}
```

### 2. Role-Based Access Control

```typescript
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);

  // Check if user has admin role (assuming role in token or separate API call)
  const response = await fetch(`${API_URL}/api/users/${user.userId}/roles`, {
    headers: { 'Authorization': `Bearer ${user.token}` },
  });

  const { roles } = await response.json();

  if (!roles.includes('admin')) {
    throw redirect('/dashboard'); // Not authorized
  }

  return { user, roles };
}
```

### 3. Combined Data Loading

```typescript
export async function loader({ request }: Route.LoaderArgs) {
  const user = await requireAuth(request);

  // Load multiple resources in parallel
  const [profileRes, settingsRes, notificationsRes] = await Promise.all([
    fetch(`${API_URL}/api/users/${user.userId}/profile`, {
      headers: { 'Authorization': `Bearer ${user.token}` },
    }),
    fetch(`${API_URL}/api/users/${user.userId}/settings`, {
      headers: { 'Authorization': `Bearer ${user.token}` },
    }),
    fetch(`${API_URL}/api/users/${user.userId}/notifications`, {
      headers: { 'Authorization': `Bearer ${user.token}` },
    }),
  ]);

  return {
    user,
    profile: await profileRes.json(),
    settings: await settingsRes.json(),
    notifications: await notificationsRes.json(),
  };
}
```

---

## Troubleshooting

### Issue: Infinite Redirect Loop

**Symptoms**: Browser shows "too many redirects" error

**Causes**:
1. Login page has `requireAuth` in loader
2. Protected page redirects to itself

**Solution**:
```typescript
// ❌ Don't do this on login page
export async function loader({ request }: Route.LoaderArgs) {
  await requireAuth(request); // Will redirect to /login if not auth
  // ...
}

// ✅ Do this instead - optional auth on login
export async function loader({ request }: Route.LoaderArgs) {
  const user = await getAuthUser(request);

  if (user) {
    // Already logged in, redirect to dashboard
    throw redirect('/dashboard');
  }

  return {};
}
```

### Issue: User Data Not Available

**Symptoms**: `user.userName` is undefined

**Cause**: Session keys don't match

**Solution**:
```typescript
// Check session keys in login action
session.set("userId", id);     // ✅ Correct
session.set("userName", name); // ✅ Correct
session.set("name", name);     // ❌ Wrong key

// In loader
const userName = session.get("userName"); // ✅ Matches action
```

### Issue: Session Expires Too Soon

**Symptoms**: User logged out unexpectedly

**Cause**: Token `exp` not set correctly

**Solution**:
```typescript
// Ensure API returns exp field
const { exp } = result.data;
const maxAge = calculateSessionMaxAge(exp);

// Check exp value
console.log('Token expires at:', new Date(exp * 1000));
console.log('Session maxAge:', maxAge, 'seconds');
```

---

## Testing

### Manual Testing

1. **Test Protected Route Access**:
```bash
# Not logged in
curl -I http://localhost:5173/dashboard
# Should redirect to /login (302)

# Logged in
curl -I http://localhost:5173/dashboard -H "Cookie: __session=..."
# Should return 200
```

2. **Test Logout**:
```bash
# POST to logout
curl -X POST http://localhost:5173/logout -H "Cookie: __session=..."
# Should destroy session and redirect
```

### Automated Testing (Future)

```typescript
describe('Dashboard Protection', () => {
  it('redirects to login if not authenticated', async () => {
    const response = await fetch('/dashboard');
    expect(response.status).toBe(302);
    expect(response.headers.get('Location')).toBe('/login');
  });

  it('renders dashboard if authenticated', async () => {
    const response = await fetch('/dashboard', {
      headers: { Cookie: `__session=${validSessionCookie}` },
    });
    expect(response.status).toBe(200);
  });
});
```

---

## Migration Checklist

If migrating from client-side auth to loader-based auth:

- [ ] Add loader to protected routes
- [ ] Replace `useSession()` with `useLoaderData()`
- [ ] Remove `useEffect` auth checks
- [ ] Remove client-side navigation logic
- [ ] Update component to use loader data
- [ ] Test protected route access
- [ ] Test logout functionality
- [ ] Update tests if any

---

## Summary

### Files Changed

| File | Status | Purpose |
|------|--------|---------|
| `routes/logout.tsx` | ✅ Created | Logout action with session destruction |
| `routes/dashboard.tsx` | ✅ Enhanced | Protected layout with loader |
| `routes/dashboard._index.tsx` | ✅ Enhanced | Protected page with loader |
| `lib/auth-guard.server.ts` | ✅ Created | `requireAuth()` and `getAuthUser()` |

### Key Improvements

1. ✅ **Server-side authentication** (secure)
2. ✅ **No flash of content** (better UX)
3. ✅ **Simpler components** (less code)
4. ✅ **Type-safe loader data** (fewer bugs)
5. ✅ **Proper logout** (session destruction)
6. ✅ **Easy to add protected routes** (just add loader)

### Next Steps

1. Apply loader pattern to other protected routes
2. Add role-based access control if needed
3. Add API calls in loaders for data fetching
4. Update tests for new auth pattern

Your dashboard is now properly protected with React Router's server-side authentication! 🔒✨
