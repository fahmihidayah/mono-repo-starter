import { deleteResourceAction } from "~/lib/actions";
import { userApi } from "../api/users";

export async function deleteUserAction({ request }: { request: Request }) {
  return deleteResourceAction({
    request,
    deleteFn: async (id, req) => {
      return userApi.deleteById({ request: req, id });
    },
    deleteManyFn: async (ids, req) => {
      return userApi.deleteMany({ request: req, ids });
    },
    resourceName: "user",
  });
}
