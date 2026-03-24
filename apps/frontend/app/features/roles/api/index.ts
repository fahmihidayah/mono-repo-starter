import type { Role, RoleFormSchema } from "~/features/roles/types";
import { ApiClient } from "../../../lib/api";

export const roleApi = new ApiClient<Role, RoleFormSchema>({
  resource: "api/roles",
  defaultToken: "",
});
