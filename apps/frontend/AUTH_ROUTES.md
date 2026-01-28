# Auth Routes Organization

The authentication routes have been organized using React Router v7's pathless layout feature for better maintainability while keeping clean URLs.

## Directory Structure

```
app/routes/
├── _auth.tsx          → Pathless layout wrapper
├── _auth/
│   ├── login.tsx      → /login route
│   └── register.tsx   → /register route
├── dashboard/
│   └── ...
└── ...
```

## How Pathless Layouts Work

The underscore prefix (`_auth`) creates a **pathless layout route**. This means:
- The folder name `_auth` is NOT included in the URL
- Files inside `_auth/` create routes at the root level
- The `_auth.tsx` layout wraps all child routes
- Perfect for grouping related routes without affecting URLs

## Routes

| File | Route Path | Description |
|------|-----------|-------------|
| `_auth/login.tsx` | `/login` | User login page |
| `_auth/register.tsx` | `/register` | User registration page |

## Updated Links

All navigation links use the clean paths:

### From Homepage ([app/routes/_index.tsx](app/routes/_index.tsx))
- Navigation "Sign in" button → `/login`
- Navigation "Get Started" button → `/register`
- Hero "Get Started Free" button → `/register`

### From Dashboard ([app/routes/dashboard.tsx](app/routes/dashboard.tsx))
- Unauthorized redirect → `/login`

### From Auth Pages
- Login page "Sign up" link → `/register`
- Register page "Sign in" link → `/login`

### From Auth Guard ([app/lib/auth-guard.tsx](app/lib/auth-guard.tsx))
- Unauthorized redirect → `/login`

## Benefits of This Organization

1. **Clean URLs**: Routes are at `/login` and `/register` (no `/auth/` prefix)
2. **Better Organization**: Auth-related files grouped in `_auth/` directory
3. **Easier Maintenance**: All authentication pages in one place
4. **Shared Layout**: Can add common auth layout/styling in `_auth.tsx`
5. **Scalability**: Easy to add more auth pages (forgot password, reset, etc.)

## The Layout File

The `_auth.tsx` file is a pathless layout that wraps all auth routes:

```tsx
import { Outlet } from "react-router";

export default function Auth() {
  return <Outlet />;
}
```

You can enhance this to add:
- Common background/styling for auth pages
- Auth-specific headers/footers
- Redirect logic for already-authenticated users
- Loading states

## Adding More Auth Routes

To add new auth routes, simply create new files in the `_auth/` directory:

```
_auth/
├── login.tsx
├── register.tsx
├── forgot-password.tsx    → /forgot-password
└── reset-password.tsx     → /reset-password
```

React Router will automatically create routes based on the file structure.

## React Router v7 Naming Conventions

- `_folder` - Creates pathless layout (folder name not in URL)
- `folder` - Creates route segment (folder name in URL)
- `file.tsx` - Creates route at `/file`
- `_file.tsx` - Creates pathless route (no URL segment)
- `file._index.tsx` - Index route for `/file`
- `file.$param.tsx` - Dynamic route for `/file/:param`

## Example URL Access

Once the dev server is running:
- Login: http://localhost:5173/login
- Register: http://localhost:5173/register

All navigation and redirects work correctly with these clean paths!
