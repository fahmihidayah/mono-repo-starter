import { apiClient } from '../api-client';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
}

// API Response structure matching your Go API
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// Auth data structure from your API
export interface AuthData {
  id: string;
  name: string;
  email: string;
  token: string;
  exp : number;
  created_at: string;
}

// User session structure (for cookies)
export interface UserSession {
  id: string;
  name: string;
  email: string;
  createdAt: string;
}

export const authApi = {
  // Login - POST /api/users/auth/login
  login: async (credentials: LoginCredentials): Promise<ApiResponse<AuthData>> => {
    return apiClient.post<ApiResponse<AuthData>>('/api/users/auth/login', credentials);
  },

  // Register - POST /api/users/auth/register
  register: async (data: RegisterData): Promise<ApiResponse<AuthData>> => {
    return apiClient.post<ApiResponse<AuthData>>('/api/users/auth/register', data);
  },

  // Logout - POST /api/users/auth/logout
  logout: async (): Promise<ApiResponse<null>> => {
    return apiClient.post<ApiResponse<null>>('/api/users/auth/logout');
  },

  // Get current user - GET /api/users/me (adjust if needed)
  getCurrentUser: async (): Promise<ApiResponse<UserSession>> => {
    return apiClient.get<ApiResponse<UserSession>>('/api/users/me');
  },
};
