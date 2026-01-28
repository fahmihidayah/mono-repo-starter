import { createCookieSessionStorage } from "react-router";

/**
 * Session data stored in the cookie
 */
export interface SessionData {
  userId: string;
  token: string;
  userName: string;
  userEmail: string;
}

/**
 * Flash data for one-time messages
 */
export interface SessionFlashData {
  error: string;
  success: string;
}

/**
 * Cookie-based session storage using React Router's official API
 * Best practices from: https://reactrouter.com/explanation/sessions-and-cookies
 */
const { getSession, commitSession, destroySession } =
  createCookieSessionStorage<SessionData, SessionFlashData>({
    cookie: {
      name: "__session",

      // Security settings
      httpOnly: true,  // Prevent client-side JavaScript access
      path: "/",
      sameSite: "lax", // CSRF protection
      secrets: [process.env.SESSION_SECRET || "dev-secret-change-in-production"],
      secure: process.env.NODE_ENV === "production", // HTTPS only in production

      // Default maxAge (7 days) - can be overridden per session
      maxAge: 60 * 60 * 24 * 7, // 7 days in seconds
    },
  });

export { getSession, commitSession, destroySession };