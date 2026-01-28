import { taskApi, type Task, type PaginatedResponse } from "~/lib/api/tasks";

/**
 * Task Repository
 * Wrapper around task API calls
 */
export class TaskRepository {
  async findById(id: string): Promise<Task | undefined> {
    try {
      return await taskApi.getById(id);
    } catch (error) {
      return undefined;
    }
  }

  async findManyPaginated(options: {
    page?: number;
    pageSize?: number;
  } = {}): Promise<PaginatedResponse<Task>> {
    const { page = 1, pageSize = 10 } = options;
    return taskApi.getAll(page, pageSize);
  }

  async create(data: Partial<Task>): Promise<Task> {
    return taskApi.create(data);
  }

  async update(id: string, data: Partial<Task>): Promise<Task | undefined> {
    try {
      return await taskApi.update(id, data);
    } catch (error) {
      return undefined;
    }
  }

  async delete(id: string): Promise<Task | undefined> {
    try {
      await taskApi.delete(id);
      return { id } as Task;
    } catch (error) {
      return undefined;
    }
  }
}

export const taskRepository = new TaskRepository();
