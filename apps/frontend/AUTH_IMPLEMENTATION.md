# Authentication Implementation

This document describes the authentication system implemented for the frontend to connect with the Go API backend.

## Overview

The authentication system uses JWT-based authentication with the following flow:
1. User logs in via `/login` or registers via `/register`
2. Backend returns a JWT token and user information
3. Token is stored in localStorage
4. Token is automatically included in all API requests via Authorization header
5. Protected routes check for valid session before rendering

## Files Structure

### Core Authentication Files

#### [app/lib/api/auth.ts](app/lib/api/auth.ts)
- **Purpose**: API endpoints for authentication
- **Exports**:
  - `authApi.login(credentials)` - Login with email/password
  - `authApi.register(data)` - Register new user
  - `authApi.logout()` - Logout user
  - `authApi.getCurrentUser()` - Get current user info
- **Types**:
  ```typescript
  interface LoginCredentials {
    email: string;
    password: string;
  }

  interface RegisterData {
    name: string;
    email: string;
    password: string;
  }

  interface AuthResponse {
    user: {
      id: string;
      name: string;
      email: string;
      emailVerified?: boolean;
      createdAt?: Date;
      updatedAt?: Date;
    };
    token: string;
  }
  ```

#### [app/lib/auth-client.ts](app/lib/auth-client.ts)
- **Purpose**: Client-side authentication utilities and session management
- **Features**:
  - Session storage in localStorage
  - React hooks for session management
  - Sign in/sign up/sign out methods
- **Exports**:
  - `authClient.signIn.email(credentials)` - Sign in
  - `authClient.signUp.email(data)` - Sign up
  - `authClient.signOut()` - Sign out
  - `useSession()` - React hook for current session

#### [app/lib/api-client.ts](app/lib/api-client.ts)
- **Purpose**: HTTP client with automatic token injection
- **Features**:
  - Automatically adds `Authorization: Bearer <token>` header
  - Better error handling with API error messages
  - Supports GET, POST, PUT, PATCH, DELETE methods

#### [app/lib/auth-guard.tsx](app/lib/auth-guard.tsx)
- **Purpose**: Component to protect routes requiring authentication
- **Usage**:
  ```tsx
  <AuthGuard>
    <ProtectedContent />
  </AuthGuard>
  ```

### Route Files

#### [app/routes/login.tsx](app/routes/login.tsx)
- Login page with email/password form
- Error handling with toast notifications
- Redirects to `/dashboard` on success

#### [app/routes/register.tsx](app/routes/register.tsx)
- Registration page with name, email, password, and confirm password fields
- Client-side validation (password matching, minimum length)
- Redirects to `/dashboard` on success

#### [app/routes/dashboard.tsx](app/routes/dashboard.tsx)
- Protected route - requires authentication
- Auto-redirects to `/login` if not authenticated
- Includes sign out functionality

## Backend API Requirements

The frontend expects the Go API to implement these endpoints:

### POST /auth/login
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "user-id",
    "name": "John Doe",
    "email": "user@example.com",
    "emailVerified": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "token": "jwt-token-here"
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "Invalid credentials",
  "message": "Invalid email or password"
}
```

### POST /auth/register
**Request Body:**
```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "user": {
    "id": "user-id",
    "name": "John Doe",
    "email": "user@example.com",
    "emailVerified": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "token": "jwt-token-here"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Validation error",
  "message": "Email already exists"
}
```

### POST /auth/logout
**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

### GET /auth/me
**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "id": "user-id",
  "name": "John Doe",
  "email": "user@example.com",
  "emailVerified": false,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

## Protected API Requests

All other API endpoints (tasks, users, etc.) should expect the `Authorization` header:

```
Authorization: Bearer <jwt-token>
```

The frontend automatically adds this header to all requests when a user is logged in.

## Session Management

- **Storage**: JWT token and user info are stored in `localStorage`
- **Keys**:
  - `auth_token` - JWT token
  - `auth_session` - User object (JSON stringified)
- **Persistence**: Session persists across page refreshes
- **Expiration**: Token expiration should be handled by the backend (401 response triggers re-login)

## Usage Examples

### Login Flow
```typescript
// User submits login form
const response = await authClient.signIn.email({
  email: 'user@example.com',
  password: 'password123'
});

if (response.error) {
  // Handle error
  toast.error(response.error.message);
} else {
  // Success - token is automatically stored
  navigate('/dashboard');
}
```

### Logout Flow
```typescript
await authClient.signOut();
// Token and session are cleared from localStorage
navigate('/');
```

### Using Session in Components
```typescript
function MyComponent() {
  const { data: session, isPending } = useSession();

  if (isPending) return <Loading />;
  if (!session) return <Login />;

  return <div>Welcome, {session.user.name}!</div>;
}
```

### Making Authenticated API Calls
```typescript
// Token is automatically included
const tasks = await taskApi.getAll();
```

## Security Considerations

1. **XSS Protection**: Token stored in localStorage is vulnerable to XSS. Ensure your app has proper XSS protections.
2. **HTTPS**: Always use HTTPS in production to protect tokens in transit.
3. **Token Expiration**: Implement token refresh or re-authentication on token expiry.
4. **CORS**: Configure your Go API to accept requests from your frontend domain.

## Testing

To test the authentication flow:

1. Start the Go API backend
2. Start the frontend: `pnpm turbo dev --filter=frontend`
3. Navigate to `/register` to create an account
4. Navigate to `/login` to sign in
5. After login, you should be redirected to `/dashboard`
6. Try accessing protected routes

## Environment Variables

Configure the API base URL in your `.env` file:

```
API_BASE_URL=http://localhost:8080
```

Or set it in your environment for production deployments.
