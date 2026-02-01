import type { BaseResponse } from '~/types';
import { ApiClient, apiClient } from '../api-client';
import type { ApiProps } from '../type';
import { getToken } from '../utils.server';

export interface User {
  id: string;
  name: string;
  email: string;
  emailVerified: boolean;
  image?: string;
  createdAt: Date;
  updatedAt: Date;
}


export const userApi = new ApiClient<User>({
  resource: "api/users",
  defaultToken: ""
});

// {
//   getAll: async ({ request, page = 1, pageSize = 10 }: ApiProps & {
//     page?: number;
//     pageSize?: number;
//   }): Promise<BaseResponse<User[]>> => {
//     const token = await getToken({ request: request });
//     return apiClient.get<User[]>(`/api/users?page=${page}&pageSize=${pageSize}`, { token });
//   },

//   getById: async ({ request, id }: ApiProps & {
//     id: string
//   }): Promise<BaseResponse<User>> => {
//     const token = await getToken({ request: request });
//     return apiClient.get<User>(`/api/users/${id}`, { token });
//   },

//   create: async ({
//     data, request
//   }: ApiProps & { data: Partial<User> }): Promise<BaseResponse<User>> => {
//     const token = await getToken({ request: request });
//     return apiClient.post<User>('/api/users', data, { token });
//   },

//   update: async ({ request, data, id }: ApiProps & { data: Partial<User>, id: string }): Promise<BaseResponse<User>> => {
//     const token = await getToken({ request: request });
//     return apiClient.put<User>(`/api/users/${id}`, data, { token });
//   },

//   delete: async ({ request, id }: ApiProps & { id: string }): Promise<BaseResponse<null>> => {
//     const token = await getToken({ request: request });
//     return apiClient.delete<null>(`/api/users/${id}`, { token });
//   },
// };
