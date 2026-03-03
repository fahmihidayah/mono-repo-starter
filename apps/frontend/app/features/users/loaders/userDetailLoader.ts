import { createDetailLoader } from "~/lib/loaders";
import { userApi } from "../api/users";

export async function userDetailLoader({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const user = await createDetailLoader({
    request,
    id,
    getByIdFn: async (userId, req) => {
      return userApi.getById({ request: req, id: userId });
    },
  });

  return { user };
}
