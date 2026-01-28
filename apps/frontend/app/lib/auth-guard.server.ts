/**
 * Server-Side Authentication Guard
 *
 * Utilities for protecting routes with authentication.
 * Use in loaders to ensure users are authenticated.
 */

import { redirect } from "react-router";
import { getSession } from "~/session.server";

export interface AuthenticatedUser {
  userId: string;
  token: string;
  userName: string;
  userEmail: string;
}

/**
 * Require authentication in a loader
 * Redirects to login if not authenticated
 *
 * @param request - The request object
 * @returns User session data
 */
export async function requireAuth(request: Request): Promise<AuthenticatedUser> {
  const session = await getSession(request.headers.get("Cookie"));

  const userId = session.get("userId");
  const token = session.get("token");
  const userName = session.get("userName") || "";
  const userEmail = session.get("userEmail") || "";

  if (!userId || !token) {
    throw redirect("/login");
  }

  return {
    userId,
    token,
    userName,
    userEmail,
  };
}

/**
 * Get current user from session (optional)
 * Returns null if not authenticated
 *
 * @param request - The request object
 * @returns User session data or null
 */
export async function getAuthUser(request: Request): Promise<AuthenticatedUser | null> {
  const session = await getSession(request.headers.get("Cookie"));

  const userId = session.get("userId");
  const token = session.get("token");
  const userName = session.get("userName") || "";
  const userEmail = session.get("userEmail") || "";

  if (!userId || !token) {
    return null;
  }

  return {
    userId,
    token,
    userName,
    userEmail,
  };
}
