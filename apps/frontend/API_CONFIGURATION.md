# API Configuration Guide

This document explains how to configure the frontend to connect to the Go API backend.

## Issue: "Route Not Found" Error

If you're seeing this error when trying to login:
```
[Login Action] Error: Error: Route Not Found at ApiClient.request
```

This means the frontend cannot connect to the API backend.

## Solution

### 1. Environment Variable Configuration

**IMPORTANT:** In Vite (used by React Router), environment variables must be prefixed with `VITE_`.

#### ✅ Correct Configuration

**File: `.env`**
```bash
VITE_API_BASE_URL=http://localhost:8080
```

**File: `app/lib/api-client.ts`**
```typescript
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
```

#### ❌ Wrong Configuration

```bash
# This will NOT work in Vite
API_BASE_URL=http://localhost:8080
```

```typescript
// This will NOT work in Vite
const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
```

### 2. Start the Go API Backend

Make sure your Go API is running on port 8080:

```bash
# In one terminal - Start the API
pnpm turbo dev --filter=api

# Or if using Make directly in the api directory
cd apps/api
make dev
```

Verify the API is running:
```bash
curl http://localhost:8080/health
# Should return a 200 OK response
```

### 3. Start the Frontend

```bash
# In another terminal - Start the frontend
pnpm turbo dev --filter=frontend
```

### 4. Restart Development Server

After changing environment variables, you must restart the dev server:

1. Stop the frontend dev server (Ctrl+C)
2. Start it again: `pnpm turbo dev --filter=frontend`

**Note:** Vite only reads `.env` files at startup, not during hot reload.

## Environment Variables

### Available Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_API_BASE_URL` | Base URL for API requests | `http://localhost:8080` | No |

### TypeScript Support

The `env.d.ts` file provides TypeScript types for environment variables:

```typescript
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
}
```

This enables IntelliSense and type checking for `import.meta.env.VITE_API_BASE_URL`.

## Testing the API Connection

### 1. Check API is Running

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test auth endpoint exists (should return 401 or 405)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test"}'
```

### 2. Check Browser Console

Open browser DevTools (F12) and check:

1. **Network Tab:**
   - Look for requests to `/auth/login`
   - Check the request URL (should be `http://localhost:8080/auth/login`)
   - Check response status and body

2. **Console Tab:**
   - Look for any CORS errors
   - Check for the API base URL being used

### 3. Debug API Client

Add temporary logging to see what URL is being called:

```typescript
// In app/lib/api-client.ts
private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${this.baseUrl}${endpoint}`;
  console.log('API Request URL:', url); // Add this line

  const response = await fetch(url, config);
  // ...
}
```

## Common Issues

### Issue 1: CORS Errors

**Symptom:**
```
Access to fetch at 'http://localhost:8080/auth/login' from origin 'http://localhost:5173'
has been blocked by CORS policy
```

**Solution:**
Configure CORS in your Go API to allow `http://localhost:5173`:

```go
// In your Go API
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

### Issue 2: Wrong Port

**Symptom:**
```
Failed to fetch
net::ERR_CONNECTION_REFUSED
```

**Solution:**
1. Check your Go API is actually running on port 8080
2. Verify the port in `.env` matches your API port
3. Check for port conflicts (is something else using 8080?)

### Issue 3: Route Not Found (404)

**Symptom:**
```
Route Not Found
```

**Solution:**
The API endpoint `/auth/login` doesn't exist in your Go backend. You need to implement it.

**Required Go API Endpoints:**

```go
// POST /auth/login
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    User  User   `json:"user"`
    Token string `json:"token"`
}

// POST /auth/register
type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// POST /auth/logout
// GET /auth/me
```

### Issue 4: Environment Variable Not Loading

**Symptom:**
API calls go to wrong URL or default `http://localhost:8080`

**Checklist:**
- [ ] Environment variable starts with `VITE_` prefix
- [ ] Using `import.meta.env.VITE_API_BASE_URL` not `process.env`
- [ ] Dev server was restarted after changing `.env`
- [ ] `.env` file is in the frontend root (not workspace root)

## Production Configuration

For production deployments:

```bash
# .env.production
VITE_API_BASE_URL=https://api.yourdomain.com
```

Build with production env:
```bash
pnpm turbo build --filter=frontend
```

The build will use `.env.production` values.

## Debugging Checklist

When login fails with API errors:

1. [ ] Go API is running (`pnpm turbo dev --filter=api`)
2. [ ] API is accessible (`curl http://localhost:8080/health`)
3. [ ] `.env` has `VITE_API_BASE_URL=http://localhost:8080`
4. [ ] Frontend dev server was restarted after env changes
5. [ ] Check browser Network tab for actual request URL
6. [ ] Check browser Console for errors
7. [ ] Verify `/auth/login` endpoint exists in Go API
8. [ ] CORS is configured correctly in Go API

## Example: Full Working Setup

### Terminal 1 - API
```bash
cd /Users/s/Documents/project-mono-repo/mono-starter
pnpm turbo dev --filter=api
# ✓ API running on http://localhost:8080
```

### Terminal 2 - Frontend
```bash
cd /Users/s/Documents/project-mono-repo/mono-starter
pnpm turbo dev --filter=frontend
# ✓ Frontend running on http://localhost:5173
```

### Browser
```
Navigate to: http://localhost:5173/login
Fill in credentials and submit
Expected: Request to http://localhost:8080/auth/login
```

## Quick Fix Script

```bash
#!/bin/bash
# fix-api-connection.sh

# 1. Update .env
cd apps/frontend
echo "VITE_API_BASE_URL=http://localhost:8080" > .env

# 2. Restart services
cd ../..
pnpm turbo dev --filter=api &
pnpm turbo dev --filter=frontend
```

Save this as `fix-api-connection.sh` and run it to quickly reset everything.
