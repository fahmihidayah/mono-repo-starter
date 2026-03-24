import { createResourceAction } from "~/lib/actions";
import { roleApi } from "../api";
import { roleFormSchema } from "../types";

export async function createRoleAction({ request }: { request: Request }) {
  return createResourceAction({
    request,
    schema: roleFormSchema,
    createFn: async (data, req) => {
      return roleApi.create({ request: req, data });
    },
    redirectPath: "/dashboard/roles",
    resourceName: "role",
  });
}
