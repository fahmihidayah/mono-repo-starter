import { deleteResourceAction } from "~/lib/actions";
import { roleApi } from "../api";

export async function deleteRoleAction({ request }: { request: Request }) {
  return deleteResourceAction({
    request,
    deleteFn: async (id, req) => {
      return roleApi.deleteById({ request: req, id });
    },
    deleteManyFn: async (ids, req) => {
      return roleApi.deleteMany({ request: req, ids });
    },
    resourceName: "role",
  });
}
