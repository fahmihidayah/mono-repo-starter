import { ApiClient } from "../../../lib/api";
import type { User, UserFormSchema } from "../types";

export const userApi = new ApiClient<User, UserFormSchema>({
  resource: "api/users",
  defaultToken: "",
});
