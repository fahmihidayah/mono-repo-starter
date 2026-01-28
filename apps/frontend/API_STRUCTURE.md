# API Structure & Cookie-Based Authentication

This document describes the API endpoint structure and cookie-based session management implementation.

## API Endpoint Structure

Your Go API uses the following structure:

```
/api/{resource}/auth/{action}
```

### Authentication Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/users/auth/login` | POST | User login |
| `/api/users/auth/register` | POST | User registration |
| `/api/users/auth/logout` | POST | User logout |
| `/api/users/me` | GET | Get current user |

### API Response Format

All endpoints return a consistent response structure:

```typescript
{
  code: number,      // HTTP status code (200, 400, etc.)
  message: string,   // Human-readable message
  data: T            // Response data (type varies by endpoint)
}
```

### Login Response Example

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "id": "1ced7b05-0cca-43ed-ae69-5f79d99821a3",
    "name": "Fah me",
    "email": "fahmi@gmail.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "created_at": "2026-01-28T02:54:28.136371Z"
  }
}
```

## Cookie-Based Session Management

### Why Cookies Instead of localStorage?

**Security Benefits:**
- ✅ Can be HttpOnly (not accessible via JavaScript)
- ✅ Automatically sent with requests
- ✅ Better CSRF protection with SameSite
- ✅ More secure against XSS attacks

**localStorage Issues:**
- ❌ Accessible via any JavaScript code
- ❌ Vulnerable to XSS attacks
- ❌ Not sent automatically with requests

### Cookie Structure

The application uses two cookies:

#### 1. `auth_token` Cookie
Stores the JWT token for API authentication.

```
auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Path=/
Max-Age=604800  (7 days)
SameSite=Lax
HttpOnly        (in production)
Secure          (in production with HTTPS)
```

#### 2. `auth_session` Cookie
Stores user session data for client-side access.

```
auth_session={"id":"...","name":"...","email":"...","createdAt":"..."}
Path=/
Max-Age=604800  (7 days)
SameSite=Lax
```

## Login Action Implementation

### File: `app/routes/_auth.login/route.tsx`

```typescript
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const email = formData.get("email") as string;
  const password = formData.get("password") as string;

  // 1. Validate input
  // ... validation code ...

  // 2. Call Go API
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
  const apiResponse = await fetch(`${apiBaseUrl}/api/users/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  // 3. Parse response
  const response = await apiResponse.json();

  // Response structure: { code: 200, message: "...", data: { token, email, name, id, created_at } }
  const { token, email: userEmail, name, id, created_at } = response.data;

  // 4. Create cookies
  const headers = new Headers();

  // Token cookie
  headers.append('Set-Cookie',
    `auth_token=${token}; Path=/; Max-Age=604800; SameSite=Lax`
  );

  // Session cookie
  const userSession = { id, name, email: userEmail, createdAt: created_at };
  headers.append('Set-Cookie',
    `auth_session=${encodeURIComponent(JSON.stringify(userSession))}; Path=/; Max-Age=604800; SameSite=Lax`
  );

  // 5. Redirect with cookies
  return redirect("/dashboard", { headers });
}
```

## Cookie Utilities

### File: `app/lib/cookies.server.ts`

Server-side cookie utilities for React Router actions/loaders:

```typescript
import { parseCookies, getAuthToken, getUserSession, createAuthCookies, clearAuthCookies, requireAuth } from '~/lib/cookies.server';

// Parse cookies from request
const cookies = parseCookies(request.headers.get('Cookie'));

// Get auth token
const token = getAuthToken(request);

// Get user session
const user = getUserSession(request);

// Create auth cookies for login
const headers = createAuthCookies(token, userSession);

// Clear cookies for logout
const headers = clearAuthCookies();

// Require auth in protected routes
const { token, user } = requireAuth(request); // Throws redirect if not authenticated
```

## Protected Route Example

### Using Cookies in Loaders

```typescript
// app/routes/dashboard.tasks/route.tsx
import { requireAuth } from '~/lib/cookies.server';

export async function loader({ request }: Route.LoaderArgs) {
  // Ensure user is authenticated
  const { token, user } = requireAuth(request);

  // Use token to fetch data from API
  const apiUrl = import.meta.env.VITE_API_BASE_URL;
  const response = await fetch(`${apiUrl}/api/tasks`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  const data = await response.json();
  return { tasks: data.data, user };
}
```

## Client-Side Session Access

### Reading Cookies on Client

While the `auth_token` should be HttpOnly in production (not accessible to JavaScript), the `auth_session` cookie can be read for displaying user info:

```typescript
// app/lib/auth-client.ts
export function getUserFromCookies(): UserSession | null {
  if (typeof document === 'undefined') return null;

  const cookies = document.cookie.split(';').reduce((acc, cookie) => {
    const [name, value] = cookie.trim().split('=');
    if (name && value) acc[name] = decodeURIComponent(value);
    return acc;
  }, {} as Record<string, string>);

  const sessionData = cookies.auth_session;
  if (!sessionData) return null;

  try {
    return JSON.parse(sessionData);
  } catch {
    return null;
  }
}
```

## Logout Implementation

### Clear Cookies on Logout

```typescript
// app/routes/logout/route.tsx
import { clearAuthCookies } from '~/lib/cookies.server';

export async function action({ request }: Route.ActionArgs) {
  // Optionally call API logout endpoint
  const token = getAuthToken(request);
  if (token) {
    const apiUrl = import.meta.env.VITE_API_BASE_URL;
    await fetch(`${apiUrl}/api/users/auth/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    }).catch(() => {}); // Ignore errors
  }

  // Clear cookies
  const headers = clearAuthCookies();

  // Redirect to home
  return redirect('/', { headers });
}
```

## Production Configuration

### Enable Security Features

For production, update cookie settings:

```typescript
// In production environment
const headers = new Headers();

headers.append('Set-Cookie', createCookie('auth_token', token, {
  httpOnly: true,  // ✅ Not accessible via JavaScript
  secure: true,    // ✅ HTTPS only
  sameSite: 'Strict', // ✅ Strict CSRF protection
  maxAge: 604800,
}));
```

### Environment Variables

```bash
# .env.production
VITE_API_BASE_URL=https://api.yourdomain.com
```

## Security Checklist

- [ ] Use `HttpOnly` for `auth_token` in production
- [ ] Use `Secure` flag with HTTPS
- [ ] Set appropriate `SameSite` policy
- [ ] Set reasonable expiry times (7 days)
- [ ] Validate tokens on every protected request
- [ ] Clear cookies on logout
- [ ] Use HTTPS in production
- [ ] Implement CSRF protection
- [ ] Add rate limiting on auth endpoints
- [ ] Monitor for suspicious activity

## Testing

### Check Cookies in Browser

1. Open DevTools (F12)
2. Go to Application tab → Cookies
3. Look for `http://localhost:5173`
4. Should see:
   - `auth_token` with JWT value
   - `auth_session` with user data

### Test Login Flow

```bash
# 1. Submit login form
# 2. Check Network tab
POST http://localhost:8080/api/users/auth/login
Response: { code: 200, message: "Login successful", data: { ... } }

# 3. Check Application tab
Cookies:
  auth_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
  auth_session: {"id":"...","name":"...","email":"..."}

# 4. Should redirect to /dashboard
```

### Test Protected Route

```bash
# 1. Access protected route
GET http://localhost:5173/dashboard

# 2. Server reads cookies
# 3. Makes API request with token
GET http://localhost:8080/api/tasks
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 4. Returns data
```

## Summary

- ✅ API endpoints: `/api/{resource}/auth/{action}`
- ✅ Response format: `{ code, message, data }`
- ✅ Login endpoint: `/api/users/auth/login`
- ✅ Token stored in `auth_token` cookie
- ✅ User data stored in `auth_session` cookie
- ✅ Cookies set via `Set-Cookie` headers
- ✅ No localStorage usage
- ✅ HttpOnly cookies in production
- ✅ Server-side cookie utilities available
- ✅ Automatic redirect on successful login
