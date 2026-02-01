import type { BaseResponse } from "~/types";
import { getToken } from "./utils.server";
import type { ApiProps } from "./type";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Import getAuthToken at runtime to avoid circular dependency
const getAuthToken = () => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('auth_token');
};

export interface RequestOptions extends RequestInit {
  token?: string; // Optional token to use for this request
}

export class ApiClient<T> {
  private baseUrl: string;
  private resource: string;
  private defaultToken?: string; // Optional default token for server-side usage


  private getResourceUrl(endpoint?: string): string {
    const base = `${this.baseUrl}/${this.resource}`;
    if (!endpoint) return base;
    if (endpoint.startsWith('?')) return `${base}${endpoint}`;
    // Remove leading slash from endpoint to avoid double slashes
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint.slice(1) : endpoint;
    return `${base}/${cleanEndpoint}`;
  }

  constructor({
    resource,
    baseUrl = API_BASE_URL,
    defaultToken
  }: {
    resource: string;
    baseUrl?: string;
    defaultToken?: string;
  }) {
    this.baseUrl = baseUrl;
    this.resource = resource;
    this.defaultToken = defaultToken;
  }

  private async request<T>(
    endpoint: string,
    options: RequestOptions = {}
  ): Promise<BaseResponse<T>> {
    const url = endpoint;

    // Get auth token from:
    // 1. Request-specific token (highest priority)
    // 2. Instance default token (for server-side)
    // 3. localStorage (for client-side)
    const token = options.token || this.defaultToken || getAuthToken();

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    // Add Authorization header if token exists
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    console.log("API Request to:", url, "with options:", headers);
    // Merge with any existing headers
    if (options.headers) {
      const existingHeaders = options.headers as Record<string, string>;
      Object.assign(headers, existingHeaders);
    }

    // Remove custom token property before passing to fetch
    const { token: _, ...fetchOptions } = options;

    const config: RequestInit = {
      ...fetchOptions,
      headers,
    };

    const response = await fetch(url, config);

    if (!response.ok) {
      const errorText = await response.text();
      let errorMessage = `API Error: ${response.status} ${response.statusText}`;
      console.error("API Error Response:", errorText);
      try {
        const errorJson = JSON.parse(errorText);
        errorMessage = errorJson.message || errorJson.error || errorMessage;
      } catch {
        // If response is not JSON, use the text or default message
        errorMessage = errorText || errorMessage;
      }
      return {
        code: response.status,
        message: errorMessage,
        data: undefined,
      };
    }

    return response.json();
  }


  async getAll({ request, page = 1, pageSize = 10 }: ApiProps & {
    page?: number;
    pageSize?: number;
  }): Promise<BaseResponse<T[]>> {
    const token = await getToken({ request: request });
    return this.request<T[]>(this.getResourceUrl(`?page=${page}&pageSize=${pageSize}`), { token });
  }

  async getById({ request, id }: ApiProps & {
    id: string
  }): Promise<BaseResponse<T>> {
    const token = await getToken({ request: request });
    return this.get<T>(this.getResourceUrl(id), { token });
  }

  async create({ request, data }: ApiProps & {
    data: Partial<T>
  }): Promise<BaseResponse<T>> {
    const token = await getToken({ request: request });
    return this.post(this.getResourceUrl(), data, { token });
  }

  async update({ request, id, data }: ApiProps & {
    id: string;
    data: Partial<T>
  }): Promise<BaseResponse<T>> {
    const token = await getToken({ request: request });
    return this.put(this.getResourceUrl(id), data, { token });
  }

  async deleteById({ request, id }: ApiProps & {
    id: string
  }): Promise<BaseResponse<T>> {
    const token = await getToken({ request: request });
    return this.delete(this.getResourceUrl(id), { token });
  }

  async get<T>(endpoint?: string, options?: RequestOptions): Promise<BaseResponse<T>> {
    return this.request<T>(this.getResourceUrl(endpoint), { ...options, method: 'GET' });
  }

  async post<T>(endpoint?: string, data?: unknown, options?: RequestOptions): Promise<BaseResponse<T>> {
    return this.request<T>(this.getResourceUrl(endpoint), {
      ...options,
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T>(endpoint?: string, data?: unknown, options?: RequestOptions): Promise<BaseResponse<T>> {
    return this.request<T>(this.getResourceUrl(endpoint), {
      ...options,
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async patch<T>(endpoint?: string, data?: unknown, options?: RequestOptions): Promise<BaseResponse<T>> {
    return this.request<T>(this.getResourceUrl(endpoint), {
      ...options,
      method: 'PATCH',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint?: string, options?: RequestOptions): Promise<BaseResponse<T>> {
    return this.request<T>(this.getResourceUrl(endpoint), { ...options, method: 'DELETE' });
  }
}

export const apiClient = new ApiClient({ resource: "" });
