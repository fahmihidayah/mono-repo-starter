import type { BaseResponse } from '~/types';
import { ApiClient, apiClient } from '../api';

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
