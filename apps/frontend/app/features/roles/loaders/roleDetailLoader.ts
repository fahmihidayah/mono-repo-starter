import { createDetailLoader } from "~/lib/loaders";
import { roleApi } from "../api";

export async function roleDetailLoader({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const role = await createDetailLoader({
    request,
    id,
    getByIdFn: async (roleId, req) => {
      return roleApi.getById({ request: req, id: roleId });
    },
  });

  return { role };
}
