import type { BaseResponse } from '~/types';
import { ApiClient, apiClient } from '../api-client';

export interface Task {
  id: string;
  title: string;
  description: string;
  createdAt: Date;
  updatedAt: Date;
}


export class TaskApi extends ApiClient<Task> {
  constructor() {
    super({
      resource: "api/tasks",
    });
  }
}

export const taskApi = new TaskApi();