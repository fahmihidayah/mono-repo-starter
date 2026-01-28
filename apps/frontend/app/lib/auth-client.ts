/**
 * Authentication Client
 *
 * This module provides authentication utilities for both server and client contexts.
 * - Server-side: Uses cookies for session storage (actions/loaders)
 * - Client-side: Reads cookies for session display (components)
 */

import { useState, useEffect } from 'react';
import type { UserSession, ApiResponse, AuthData } from './api/auth';

// ============================================================================
// Cookie Utilities (Works in both server and client contexts)
// ============================================================================

/**
 * Parse cookies from cookie string
 */
function parseCookies(cookieString: string | null): Record<string, string> {
  if (!cookieString) return {};

  return cookieString.split(';').reduce((cookies, cookie) => {
    const [name, value] = cookie.trim().split('=');
    if (name && value) {
      cookies[name] = decodeURIComponent(value);
    }
    return cookies;
  }, {} as Record<string, string>);
}

/**
 * Get auth token from cookies (client-side)
 */
export function getAuthToken(): string | null {
  if (typeof document === 'undefined') return null;
  const cookies = parseCookies(document.cookie);
  return cookies.auth_token || null;
}

/**
 * Get user session from cookies (client-side)
 */
export function getUserSession(): UserSession | null {
  if (typeof document === 'undefined') return null;

  const cookies = parseCookies(document.cookie);
  const sessionData = cookies.auth_session;

  if (!sessionData) return null;

  try {
    return JSON.parse(sessionData) as UserSession;
  } catch {
    return null;
  }
}

/**
 * Clear cookies (client-side) - Sets them to expire immediately
 */
export function clearCookies(): void {
  if (typeof document === 'undefined') return;

  document.cookie = 'auth_token=; Path=/; Max-Age=0';
  document.cookie = 'auth_session=; Path=/; Max-Age=0';
}

// ============================================================================
// Server-Side Authentication (for React Router actions/loaders)
// ============================================================================

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface LoginResult {
  success: boolean;
  headers?: Headers;
  error?: string;
}

/**
 * Server-side login function
 * Makes API call and returns cookies in headers for redirect
 *
 * @param credentials - Email and password
 * @returns Headers with Set-Cookie or error message
 */
export async function serverLogin(credentials: LoginCredentials): Promise<LoginResult> {
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

    // Call the Go API
    const apiResponse = await fetch(`${apiBaseUrl}/api/users/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    });

    if (!apiResponse.ok) {
      const errorData = await apiResponse.json().catch(() => ({}));
      return {
        success: false,
        error: errorData.message || errorData.error || 'Login failed',
      };
    }

    // Parse the response: { code, message, data: { token, email, name, id, created_at } }
    const response: ApiResponse<AuthData> = await apiResponse.json();

    if (response.code !== 200 || !response.data?.token) {
      return {
        success: false,
        error: response.message || 'Invalid response from server',
      };
    }

    const { token, email, name, id, created_at } = response.data;

    // Create cookies
    const headers = new Headers();

    // Set auth token cookie
    const tokenCookie = [
      `auth_token=${token}`,
      'Path=/',
      'Max-Age=604800', // 7 days
      'SameSite=Lax',
      // In production: add 'Secure' and 'HttpOnly'
    ];
    headers.append('Set-Cookie', tokenCookie.join('; '));

    // Set user session cookie
    const userSession: UserSession = {
      id,
      name,
      email,
      createdAt: created_at,
    };

    const sessionCookie = [
      `auth_session=${encodeURIComponent(JSON.stringify(userSession))}`,
      'Path=/',
      'Max-Age=604800', // 7 days
      'SameSite=Lax',
    ];
    headers.append('Set-Cookie', sessionCookie.join('; '));

    return {
      success: true,
      headers,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message || 'Network error. Please try again.',
    };
  }
}

/**
 * Server-side logout function
 * Returns headers to clear cookies
 */
export async function serverLogout(token?: string): Promise<Headers> {
  // Optionally call API logout endpoint
  if (token) {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
    await fetch(`${apiBaseUrl}/api/users/auth/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    }).catch(() => {}); // Ignore errors
  }

  // Clear cookies
  const headers = new Headers();
  headers.append('Set-Cookie', 'auth_token=; Path=/; Max-Age=0');
  headers.append('Set-Cookie', 'auth_session=; Path=/; Max-Age=0');

  return headers;
}

/**
 * Get auth data from request cookies (for loaders)
 */
export function getAuthFromRequest(request: Request): { token: string; user: UserSession } | null {
  const cookieHeader = request.headers.get('Cookie');
  const cookies = parseCookies(cookieHeader);

  const token = cookies.auth_token;
  const sessionData = cookies.auth_session;

  if (!token || !sessionData) return null;

  try {
    const user = JSON.parse(sessionData) as UserSession;
    return { token, user };
  } catch {
    return null;
  }
}

/**
 * Require authentication in loader/action
 * Throws redirect to login if not authenticated
 */
export function requireAuth(request: Request): { token: string; user: UserSession } {
  const auth = getAuthFromRequest(request);

  if (!auth) {
    throw new Response(null, {
      status: 302,
      headers: {
        Location: '/login',
      },
    });
  }

  return auth;
}

// ============================================================================
// Client-Side Hooks (for React components)
// ============================================================================

/**
 * React hook to access current user session
 * Reads from cookies (client-side)
 */
export function useSession() {
  const [session, setSession] = useState<UserSession | null>(null);
  const [isPending, setIsPending] = useState(true);

  useEffect(() => {
    // Read session from cookies
    const user = getUserSession();
    setSession(user);
    setIsPending(false);
  }, []);

  return {
    data: session ? { user: session } : null,
    isPending,
  };
}

/**
 * Client-side sign out
 * Clears cookies and optionally calls logout endpoint
 */
export async function signOut() {
  const token = getAuthToken();

  // Call logout endpoint
  if (token) {
    try {
      const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
      await fetch(`${apiBaseUrl}/api/users/auth/logout`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
    } catch (error) {
      console.error('Logout error:', error);
    }
  }

  // Clear cookies
  clearCookies();

  // Redirect to home
  window.location.href = '/';
}

// ============================================================================
// Validation Utilities
// ============================================================================

export function validateEmail(email: string): string | null {
  if (!email) return 'Email is required';
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return 'Invalid email format';
  return null;
}

export function validatePassword(password: string): string | null {
  if (!password) return 'Password is required';
  if (password.length < 6) return 'Password must be at least 6 characters';
  return null;
}
