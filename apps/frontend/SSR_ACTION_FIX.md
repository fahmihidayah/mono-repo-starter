# SSR Action Fix - Login Route

## The Problem

When submitting the login form, the request was going to `http://localhost:5173/login.data` instead of your Go API at `http://localhost:8080/auth/login`.

### Why This Happened

React Router v7 runs `action` functions on the **server-side** (SSR/SSG), not in the browser. When you use client-side utilities like `authApi.login()` in an action, they try to make requests relative to the React Router server, not your external API.

**Flow:**
```
User submits form
      ↓
React Router action runs on Node.js server (port 5173)
      ↓
authApi.login() tries to call API
      ↓
API client thinks it's running on localhost:5173
      ↓
Request goes to http://localhost:5173/auth/login ❌
```

## The Solution

Make a **direct `fetch` call** to the Go API in the action, bypassing the client-side API utilities.

### Before (Broken)

```typescript
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const email = formData.get("email") as string;
  const password = formData.get("password") as string;

  // ❌ This runs on the server and makes wrong request
  const response = await authApi.login({ email, password });

  return redirect("/dashboard");
}
```

**Problem:** `authApi` uses `apiClient` which uses `fetch()` with relative URLs, and when running server-side, it resolves to the React Router server.

### After (Fixed)

```typescript
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const email = formData.get("email") as string;
  const password = formData.get("password") as string;

  // ✅ Direct fetch to the Go API
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

  const apiResponse = await fetch(`${apiBaseUrl}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  if (!apiResponse.ok) {
    const errorData = await apiResponse.json().catch(() => ({}));
    throw new Error(errorData.message || 'Login failed');
  }

  const response = await apiResponse.json();

  // Set cookies and redirect
  const headers = new Headers();
  headers.append("Set-Cookie", `auth_token=${response.token}; Path=/; Max-Age=604800`);

  return redirect("/dashboard", { headers });
}
```

**Solution:** Use `fetch()` with an absolute URL to the Go API server.

## Key Changes

1. **Removed import**: `import { authApi } from "~/lib/api/auth";`
2. **Direct fetch**: Use `fetch()` with absolute URL
3. **Environment variable**: Get API URL from `import.meta.env.VITE_API_BASE_URL`
4. **Error handling**: Parse error response from API

## Understanding React Router SSR

### Server-Side vs Client-Side

| Context | Where Code Runs | `fetch()` behavior |
|---------|----------------|-------------------|
| `action` | Node.js server (SSR) | Relative URLs resolve to React Router server |
| `loader` | Node.js server (SSR) | Relative URLs resolve to React Router server |
| Component | Browser | Relative URLs resolve to browser URL |
| Client code | Browser | Can use relative URLs normally |

### When to Use Direct Fetch

Use direct `fetch()` with absolute URLs in:
- ✅ `action` functions
- ✅ `loader` functions
- ✅ Server-side utilities

Use `apiClient` in:
- ✅ Client components (browser only)
- ✅ Client-side event handlers
- ✅ React effects

## Network Request Flow

### Before Fix
```
Browser Form Submit
      ↓
POST http://localhost:5173/login (React Router)
      ↓
Server Action Runs
      ↓
authApi.login() called
      ↓
fetch('/auth/login') with relative URL
      ↓
Resolves to http://localhost:5173/auth/login ❌
      ↓
404 Not Found
```

### After Fix
```
Browser Form Submit
      ↓
POST http://localhost:5173/login (React Router)
      ↓
Server Action Runs
      ↓
fetch('http://localhost:8080/auth/login') with absolute URL
      ↓
Request goes to Go API ✅
      ↓
200 OK with token
      ↓
Set cookies, redirect to /dashboard
```

## Testing

### Check Network Tab

After the fix, you should see:

**React Router Request (Normal):**
```
POST http://localhost:5173/login
Type: document
Initiator: Form submission
```

**API Request (From Action):**
```
POST http://localhost:8080/auth/login
Type: fetch
Initiator: Server-side action
Status: 200 OK
Response: { user: {...}, token: "..." }
```

### Verify in Browser DevTools

1. Open Network tab (F12)
2. Submit login form
3. Should see request to `http://localhost:8080/auth/login`
4. Response should have user data and token
5. Should redirect to `/dashboard`

## Common Patterns

### Pattern 1: Login/Register Actions

```typescript
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();

  // ✅ Direct fetch to external API
  const apiUrl = import.meta.env.VITE_API_BASE_URL;
  const response = await fetch(`${apiUrl}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(Object.fromEntries(formData)),
  });

  const data = await response.json();

  // Set cookies
  return redirect("/dashboard", {
    headers: {
      "Set-Cookie": `token=${data.token}; Path=/; Max-Age=604800`
    }
  });
}
```

### Pattern 2: Data Fetching Loaders

```typescript
export async function loader({ params }: Route.LoaderArgs) {
  // ✅ Direct fetch to external API
  const apiUrl = import.meta.env.VITE_API_BASE_URL;
  const response = await fetch(`${apiUrl}/tasks/${params.id}`);

  if (!response.ok) {
    throw new Response("Not Found", { status: 404 });
  }

  return response.json();
}
```

### Pattern 3: Authenticated Requests

```typescript
export async function loader({ request }: Route.LoaderArgs) {
  // Get token from cookies
  const cookie = request.headers.get("Cookie");
  const token = cookie?.split(';')
    .find(c => c.trim().startsWith('auth_token='))
    ?.split('=')[1];

  // ✅ Direct fetch with auth
  const apiUrl = import.meta.env.VITE_API_BASE_URL;
  const response = await fetch(`${apiUrl}/api/protected`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  return response.json();
}
```

## Best Practices

### 1. Create a Server-Side API Client

Create a separate utility for server-side API calls:

```typescript
// app/lib/server-api.ts
export async function serverFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const apiUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

  const response = await fetch(`${apiUrl}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.message || `API Error: ${response.status}`);
  }

  return response.json();
}
```

Then use it in actions:

```typescript
import { serverFetch } from "~/lib/server-api";

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();

  const data = await serverFetch('/auth/login', {
    method: 'POST',
    body: JSON.stringify(Object.fromEntries(formData)),
  });

  return redirect("/dashboard");
}
```

### 2. Environment Variables

Always use environment variables for API URLs:

```typescript
// ✅ Good - configurable
const apiUrl = import.meta.env.VITE_API_BASE_URL;

// ❌ Bad - hardcoded
const apiUrl = 'http://localhost:8080';
```

### 3. Error Handling

Always handle API errors properly:

```typescript
try {
  const response = await fetch(`${apiUrl}/auth/login`, options);

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    return {
      errors: { general: error.message || 'Login failed' }
    };
  }

  return await response.json();
} catch (error) {
  return {
    errors: { general: 'Network error. Please try again.' }
  };
}
```

## Debugging Tips

### Check Where Code Runs

Add logging to see where code executes:

```typescript
export async function action({ request }: Route.ActionArgs) {
  console.log('Action running on:', typeof window === 'undefined' ? 'SERVER' : 'CLIENT');
  // Will log: "Action running on: SERVER"
}
```

### Inspect Request URL

Log the full URL being requested:

```typescript
const url = `${apiUrl}/auth/login`;
console.log('Fetching:', url);
// Should log: "Fetching: http://localhost:8080/auth/login"
```

### Check Response

Log the response to debug API issues:

```typescript
const response = await fetch(url, options);
console.log('Response status:', response.status);
console.log('Response headers:', Object.fromEntries(response.headers));
const data = await response.json();
console.log('Response data:', data);
```

## Summary

- ✅ React Router actions run **server-side**
- ✅ Use **direct `fetch()`** with absolute URLs
- ✅ Don't use client-side API utilities in actions
- ✅ Use environment variables for API URLs
- ✅ Handle errors from API responses
- ✅ Set cookies via response headers

The login now correctly calls your Go API at `http://localhost:8080/auth/login`! 🚀
