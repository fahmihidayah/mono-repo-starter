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
  logout: async (token: string): Promise<ApiResponse<null>> => {
    return apiClient.post<ApiResponse<null>>('/api/users/auth/logout', {}, {
      token,
    });
  },

  // Get current user - GET /api/users/me (adjust if needed)
  getCurrentUser: async (): Promise<ApiResponse<UserSession>> => {
    return apiClient.get<ApiResponse<UserSession>>('/api/users/me');
  },

  // Update user profile - PUT /api/users/auth/me
  updateProfile: async (
    data: { name: string; email: string },
    token?: string
  ): Promise<ApiResponse<UserSession>> => {
    return apiClient.put<ApiResponse<UserSession>>('/api/users/auth/me', data, { token });
  },

  forgotPassword: async (email: string): Promise<ApiResponse<null>> => {
    return apiClient.post<ApiResponse<null>>('/api/users/auth/initial-reset-password', { email });
  },

  // Change password - POST /api/users/auth/change-password
  changePassword: async (
    data: { currentPassword: string; newPassword: string },
    token?: string
  ): Promise<ApiResponse<null>> => {
    // Transform camelCase to snake_case for backend
    const requestData = {
      old_password: data.currentPassword,
      new_password: data.newPassword,
    };
    return apiClient.put<ApiResponse<null>>('/api/users/auth/change-password', requestData, { token });
  },
};
