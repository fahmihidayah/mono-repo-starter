import { ApiClient } from "../../../lib/api";

export interface User {
  id: string;
  name: string;
  email: string;
  password?: string;
  emailVerified: boolean;
  image?: string;
  createdAt: Date;
  updatedAt: Date;
}

export const userApi = new ApiClient<User, any>({
  resource: "api/users",
  defaultToken: "",
});
