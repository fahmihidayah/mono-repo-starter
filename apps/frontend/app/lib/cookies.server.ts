/**
 * Server-side cookie utilities for React Router actions/loaders
 * These functions only work in server-side code (actions/loaders)
 */

import type { UserSession } from './api/auth';

/**
 * Parse cookies from request header
 */
export function parseCookies(cookieHeader: string | null): Record<string, string> {
  if (!cookieHeader) return {};

  return cookieHeader.split(';').reduce((cookies, cookie) => {
    const [name, value] = cookie.trim().split('=');
    if (name && value) {
      cookies[name] = decodeURIComponent(value);
    }
    return cookies;
  }, {} as Record<string, string>);
}

/**
 * Get auth token from request cookies
 */
export function getAuthToken(request: Request): string | null {
  const cookieHeader = request.headers.get('Cookie');
  const cookies = parseCookies(cookieHeader);
  return cookies.auth_token || null;
}

/**
 * Get user session from request cookies
 */
export function getUserSession(request: Request): UserSession | null {
  const cookieHeader = request.headers.get('Cookie');
  const cookies = parseCookies(cookieHeader);

  const sessionData = cookies.auth_session;
  if (!sessionData) return null;

  try {
    return JSON.parse(sessionData) as UserSession;
  } catch {
    return null;
  }
}

/**
 * Create cookie string for Set-Cookie header
 */
export function createCookie(
  name: string,
  value: string,
  options: {
    maxAge?: number;
    path?: string;
    sameSite?: 'Lax' | 'Strict' | 'None';
    secure?: boolean;
    httpOnly?: boolean;
  } = {}
): string {
  const {
    maxAge = 604800, // 7 days default
    path = '/',
    sameSite = 'Lax',
    secure = false,
    httpOnly = false,
  } = options;

  const parts = [
    `${name}=${encodeURIComponent(value)}`,
    `Path=${path}`,
    `Max-Age=${maxAge}`,
    `SameSite=${sameSite}`,
  ];

  if (secure) parts.push('Secure');
  if (httpOnly) parts.push('HttpOnly');

  return parts.join('; ');
}

/**
 * Create auth cookies for login response
 */
export function createAuthCookies(token: string, user: UserSession): Headers {
  const headers = new Headers();

  // Token cookie (make HttpOnly in production for security)
  headers.append('Set-Cookie', createCookie('auth_token', token, {
    httpOnly: false, // Set to true in production
    secure: false,   // Set to true in production (HTTPS only)
  }));

  // Session cookie
  headers.append('Set-Cookie', createCookie('auth_session', JSON.stringify(user)));

  return headers;
}

/**
 * Clear auth cookies (for logout)
 */
export function clearAuthCookies(): Headers {
  const headers = new Headers();

  // Clear token cookie
  headers.append('Set-Cookie', createCookie('auth_token', '', { maxAge: 0 }));

  // Clear session cookie
  headers.append('Set-Cookie', createCookie('auth_session', '', { maxAge: 0 }));

  return headers;
}

/**
 * Require authenticated user in loader/action
 * Throws redirect to login if not authenticated
 */
export function requireAuth(request: Request): { token: string; user: UserSession } {
  const token = getAuthToken(request);
  const user = getUserSession(request);

  if (!token || !user) {
    throw new Response(null, {
      status: 302,
      headers: {
        Location: '/login',
      },
    });
  }

  return { token, user };
}
