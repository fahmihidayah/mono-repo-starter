import { apiClient } from '../api-client';

export interface Task {
  id: string;
  title: string;
  description: string;
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

export const taskApi = {
  getAll: async (page = 1, pageSize = 10): Promise<PaginatedResponse<Task>> => {
    return apiClient.get<PaginatedResponse<Task>>(`/tasks?page=${page}&pageSize=${pageSize}`);
  },

  getById: async (id: string): Promise<Task> => {
    return apiClient.get<Task>(`/tasks/${id}`);
  },

  create: async (data: Partial<Task>): Promise<Task> => {
    return apiClient.post<Task>('/tasks', data);
  },

  update: async (id: string, data: Partial<Task>): Promise<Task> => {
    return apiClient.put<Task>(`/tasks/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    return apiClient.delete<void>(`/tasks/${id}`);
  },
};
