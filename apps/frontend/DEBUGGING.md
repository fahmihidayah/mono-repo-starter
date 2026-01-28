# Debugging Guide

This app includes comprehensive logging to help you debug issues quickly. Open your browser's developer console to see detailed logs.

## Logging System

The app uses a custom logger utility (`app/utils/logger.ts`) that provides structured logging throughout the application.

### Log Prefixes

All logs are prefixed with the component name for easy filtering:

- `[Auth]` - Authentication-related logs (auth client, sign in/up/out)
- `[Login]` - Login page logs
- `[Register]` - Registration page logs
- `[Dashboard]` - Dashboard page logs
- `[ReactQuery]` - React Query cache and mutation logs

### Viewing Logs

1. Open your browser's Developer Tools (F12 or Cmd+Option+I on Mac)
2. Go to the Console tab
3. Use the filter to search for specific components:
   - Type `[Auth]` to see only authentication logs
   - Type `[Login]` to see login flow logs
   - Type `error` to see all errors

### What Gets Logged

#### Authentication Flow

**Login (`/login`):**
```
[Login] Form submitted { email: "..." }
[Login] Attempting sign in...
[Auth] Sign in attempt { email: "..." }
[Auth] Sign in successful { email: "..." }
[Login] Sign in response: { ... }
[Login] Sign in successful, redirecting to dashboard
```

**Registration (`/register`):**
```
[Register] Form submitted { name: "...", email: "...", passwordLength: 8 }
[Register] Validation passed, starting registration...
[Register] Attempting sign up...
[Auth] Sign up attempt { name: "...", email: "..." }
[Auth] Sign up successful { email: "..." }
[Register] Sign up response: { ... }
[Register] Attempting automatic sign in...
[Auth] Sign in attempt { email: "..." }
[Register] Auto sign in successful, redirecting to dashboard
```

**Dashboard (`/dashboard`):**
```
[Dashboard] Component mounted/updated {
  isPending: false,
  hasSession: true,
  hasError: false,
  sessionData: { userId: "...", email: "...", name: "..." }
}
[Dashboard] Auth state changed { ... }
```

**Sign Out:**
```
[Dashboard] Sign out initiated
[Dashboard] Calling signOut...
[Auth] Sign out attempt
[Auth] Sign out successful
[Dashboard] Sign out successful
```

#### React Query

The React Query provider logs all query and mutation events:

```
[ReactQuery] Creating QueryClient...
[ReactQuery] Mutation error { message: "...", error: { ... } }
[ReactQuery] Mutation success
[ReactQuery] Query retry check { failureCount: 1, errorStatus: 500 }
```

## Using the Logger in Your Code

### Import the logger

```typescript
import { logger } from '~/utils/logger';

// Or use pre-configured scoped loggers
import { authLogger, apiLogger, routeLogger } from '~/utils/logger';
```

### Basic Usage

```typescript
// Info log
logger.info('ComponentName', 'Action completed', { data: value });

// Warning log
logger.warn('ComponentName', 'Something unexpected', { details });

// Error log
logger.error('ComponentName', 'Operation failed', error);

// Debug log (only in development)
logger.debug('ComponentName', 'Debug info', { state });
```

### Create a Scoped Logger

```typescript
const myLogger = logger.scope('MyComponent');

myLogger.info('Component mounted');
myLogger.warn('Deprecated prop used', { prop: 'oldProp' });
myLogger.error('Failed to load data', error);
```

### API Request/Response Logging

```typescript
// Log API request
logger.apiRequest('POST', '/api/users', { name: 'John' });

// Log API response
logger.apiResponse('POST', '/api/users', 201, { id: '123' });
```

### Disable Logging

In production or for specific scenarios, you can disable logging:

```typescript
import { logger } from '~/utils/logger';

logger.setEnabled(false); // Disable all logs
```

## Common Debugging Scenarios

### Issue: Login not working

1. Check console for `[Login]` and `[Auth]` logs
2. Look for error messages in the response
3. Verify the auth API endpoint is responding at `/api/auth/*`
4. Check the database has the user with correct credentials

### Issue: Registration fails

1. Check console for `[Register]` logs
2. Verify password validation (minimum 8 characters)
3. Check if email already exists in database
4. Look for database connection errors

### Issue: Dashboard shows "Not authenticated"

1. Check console for `[Dashboard]` logs
2. Verify session state: `hasSession: true/false`
3. Check if cookies are being set correctly
4. Try logging in again

### Issue: Session not persisting

1. Check browser cookies for session tokens
2. Verify `BETTER_AUTH_SECRET` is set in `.env`
3. Check `BETTER_AUTH_URL` matches your dev server URL
4. Look for CORS errors in console

## Development Tools

### React Query DevTools

The app includes React Query DevTools (bottom-right corner in development):
- View all queries and their states
- See cached data
- Manually trigger refetches
- Inspect query keys and data

### Browser DevTools

**Application Tab:**
- Check cookies for auth tokens
- Inspect localStorage/sessionStorage
- View IndexedDB if used

**Network Tab:**
- Monitor API requests to `/api/auth/*`
- Check request/response payloads
- Verify status codes (200, 401, etc.)

## Tips

1. **Filter by log level**: Use console's level filter to show only warnings/errors
2. **Timestamps**: All logs include implicit browser timestamps
3. **Copy logs**: Right-click → "Save as" to export console logs
4. **Breakpoints**: Add `debugger;` statements in code for interactive debugging
5. **Network throttling**: Test slow connections in DevTools Network tab

## Getting Help

If you're stuck:
1. Check all relevant log prefixes
2. Look for error stack traces
3. Verify environment variables in `.env`
4. Run `pnpm run db:push` to ensure database schema is up to date
5. Clear browser cache and cookies, then try again
