/**
 * Client-Side Authentication Utilities
 *
 * Minimal client-side utilities for authentication.
 * Server-side authentication is handled by:
 * - ~/lib/auth.server.ts (API calls)
 * - ~/session.server.ts (Session storage with React Router)
 */

/**
 * Client-side sign out
 * Redirects to logout action which handles session destruction
 */
export async function signOut() {
  window.location.href = '/logout';
}
