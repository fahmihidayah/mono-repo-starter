import { userApi, type User, type PaginatedResponse } from "~/lib/api/users";

/**
 * User Repository
 * Wrapper around user API calls
 */
export class UserRepository {
  async findById(id: string): Promise<User | undefined> {
    try {
      return await userApi.getById(id);
    } catch (error) {
      return undefined;
    }
  }

  async findManyPaginated(options: {
    page?: number;
    pageSize?: number;
  } = {}): Promise<PaginatedResponse<User>> {
    const { page = 1, pageSize = 10 } = options;
    return userApi.getAll(page, pageSize);
  }

  async create(data: Partial<User>): Promise<User> {
    return userApi.create(data);
  }

  async update(id: string, data: Partial<User>): Promise<User | undefined> {
    try {
      return await userApi.update(id, data);
    } catch (error) {
      return undefined;
    }
  }

  async delete(id: string): Promise<User | undefined> {
    try {
      await userApi.delete(id);
      return { id } as User;
    } catch (error) {
      return undefined;
    }
  }
}

export const userRepository = new UserRepository();
