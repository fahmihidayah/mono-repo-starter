import { apiClient } from '../api-client';

export interface User {
  id: string;
  name: string;
  email: string;
  emailVerified: boolean;
  image?: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    currentPage: number;
    pageSize: number;
    totalItems: number;
    totalPages: number;
    hasNextPage: boolean;
    hasPrevPage: boolean;
    nextPage: number | null;
    prevPage: number | null;
  };
}

export const userApi = {
  getAll: async (page = 1, pageSize = 10): Promise<PaginatedResponse<User>> => {
    return apiClient.get<PaginatedResponse<User>>(`/users?page=${page}&pageSize=${pageSize}`);
  },

  getById: async (id: string): Promise<User> => {
    return apiClient.get<User>(`/users/${id}`);
  },

  create: async (data: Partial<User>): Promise<User> => {
    return apiClient.post<User>('/users', data);
  },

  update: async (id: string, data: Partial<User>): Promise<User> => {
    return apiClient.put<User>(`/users/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(`/users/${id}`);
  },
};
