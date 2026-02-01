import type { BaseResponse } from '~/types';
import { ApiClient, apiClient } from '../api-client';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
}

// Auth data structure from your API
export interface AuthData {
  id: string;
  name: string;
  email: string;
  token: string;
  exp: number;
  created_at: string;
}

// User session structure (for cookies)
export interface UserSession {
  id: string;
  name: string;
  email: string;
  createdAt: string;
}


//  {
//   // Login - POST /api/users/auth/login
//   login: async (credentials: LoginCredentials): Promise<BaseResponse<AuthData>> => {
//     return apiClient.post<AuthData>('/api/users/auth/login', credentials);
//   },

//   // Register - POST /api/users/auth/register
//   register: async (data: RegisterData): Promise<BaseResponse<AuthData>> => {
//     return apiClient.post<AuthData>('/api/users/auth/register', data);
//   },

//   // Logout - POST /api/users/auth/logout
//   logout: async (token: string): Promise<BaseResponse<null>> => {
//     return apiClient.post<null>('/api/users/auth/logout', {}, {
//       token,
//     });
//   },

//   // Get current user - GET /api/users/me (adjust if needed)
//   getCurrentUser: async (): Promise<BaseResponse<UserSession>> => {
//     return apiClient.get<UserSession>('/api/users/me');
//   },

//   // Update user profile - PUT /api/users/auth/me
//   updateProfile: async (
//     data: { name: string; email: string },
//     token?: string
//   ): Promise<BaseResponse<UserSession>> => {
//     return apiClient.put<UserSession>('/api/users/auth/me', data, { token });
//   },

//   forgotPassword: async (email: string): Promise<BaseResponse<null>> => {
//     return apiClient.post<null>('/api/users/auth/initial-reset-password', { email });
//   },

//   // Change password - POST /api/users/auth/change-password
//   changePassword: async (
//     data: { currentPassword: string; newPassword: string },
//     token?: string
//   ): Promise<BaseResponse<null>> => {
//     // Transform camelCase to snake_case for backend
//     const requestData = {
//       old_password: data.currentPassword,
//       new_password: data.newPassword,
//     };
//     return apiClient.put<null>('/api/users/auth/change-password', requestData, { token });
//   },
// };

export class AuthApi extends ApiClient<AuthData> {
  constructor() {
    super({
      resource: "api/users/auth",
    });
  }

  async login(data: LoginCredentials): Promise<BaseResponse<AuthData>> {
    return this.post("/login", data);
  }

  async register(data: RegisterData): Promise<BaseResponse<AuthData>> {
    return this.post("/register", data);
  }

  async logout(token: string): Promise<BaseResponse<null>> {
    return this.post("/logout", {}, { token });
  }

  async getCurrentUser(): Promise<BaseResponse<UserSession>> {
    return this.get("/me");
  }

  async updateProfile(
    data: { name: string; email: string },
    token?: string
  ): Promise<BaseResponse<UserSession>> {
    return this.put("/me", data, { token });
  }

  async forgotPassword(email: string): Promise<BaseResponse<null>> {
    return this.post("/initial-reset-password", { email });
  }

  async changePassword(
    data: { currentPassword: string; newPassword: string },
    token?: string
  ): Promise<BaseResponse<null>> {
    // Transform camelCase to snake_case for backend
    const requestData = {
      old_password: data.currentPassword,
      new_password: data.newPassword,
    };
    return this.put("/change-password", requestData, { token });
  }
}

export const authApi = new AuthApi();
