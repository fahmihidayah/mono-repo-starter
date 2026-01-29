/**
 * Server-side Authentication Utilities
 *
 * Handles authentication with the Go API backend.
 * Uses React Router session storage for cookie management.
 */

import type { ApiResponse, AuthData, LoginCredentials, RegisterData } from './api/auth';

/**
 * Authentication result returned by API
 */
export interface AuthResult {
  success: boolean;
  data?: AuthData;
  error?: string;
}

/**
 * Call the Go API to authenticate user
 * Does NOT set cookies - that's handled by the session storage in the action
 *
 * @param credentials - Email and password
 * @returns Authentication data or error
 */
export async function authenticateUser(
  credentials: LoginCredentials
): Promise<AuthResult> {
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

    const response = await fetch(`${apiBaseUrl}/api/users/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        error: errorData.message || 'Login failed',
      };
    }

    const apiResponse: ApiResponse<AuthData> = await response.json();

    if (apiResponse.code !== 200 || !apiResponse.data?.token) {
      return {
        success: false,
        error: apiResponse.message || 'Invalid response from server',
      };
    }

    return {
      success: true,
      data: apiResponse.data,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message || 'Network error. Please try again.',
    };
  }
}

/**
 * Call the Go API to register new user
 * Does NOT set cookies - that's handled by the session storage in the action
 *
 * @param data - Name, email, and password
 * @returns Authentication data or error
 */
export async function registerUser(
  data: RegisterData
): Promise<AuthResult> {
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

    const response = await fetch(`${apiBaseUrl}/api/users/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        error: errorData.message || 'Registration failed',
      };
    }

    const apiResponse: ApiResponse<AuthData> = await response.json();

    if (apiResponse.code !== 200 && apiResponse.code !== 201) {
      return {
        success: false,
        error: apiResponse.message || 'Invalid response from server',
      };
    }

    if (!apiResponse.data?.token) {
      return {
        success: false,
        error: 'Registration succeeded but no token received',
      };
    }

    return {
      success: true,
      data: apiResponse.data,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message || 'Network error. Please try again.',
    };
  }
}

/**
 * Calculate session maxAge from token expiration
 * @param exp - Unix timestamp of token expiration
 * @returns maxAge in seconds, or default 7 days
 */
export function calculateSessionMaxAge(exp?: number): number {
  if (!exp) {
    // Default to 7 days
    return 60 * 60 * 24 * 7;
  }

  const now = Math.floor(Date.now() / 1000); // Current time in seconds
  const maxAge = exp - now;

  // Ensure positive maxAge, minimum 60 seconds
  return Math.max(maxAge, 60);
}

/**
 * Call the Go API to logout user
 * @param token - Authentication token
 */
export async function callLogoutApi(token: string): Promise<void> {
  try {
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

    await fetch(`${apiBaseUrl}/api/users/auth/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    });
  } catch (error) {
    // Ignore logout errors - session will still be destroyed locally
    console.error('Logout API error:', error);
  }
}
