import { updateResourceAction } from "~/lib/actions";
import { roleApi } from "../api";
import { roleFormSchema } from "../types";

export async function updateRoleAction({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  return updateResourceAction({
    request,
    id,
    schema: roleFormSchema,
    updateFn: async (roleId, data, req) => {
      return roleApi.update({ request: req, data, id: roleId });
    },
    redirectPath: "/dashboard/roles",
    resourceName: "role",
  });
}
