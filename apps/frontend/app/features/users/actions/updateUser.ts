import { updateResourceAction } from "~/lib/actions";
import { userApi } from "../api/users";
import { userFormSchema } from "../types";

export async function updateUserAction({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  return updateResourceAction({
    request,
    id,
    schema: userFormSchema,
    updateFn: async (userId, data, req) => {
      return userApi.update({ request: req, data, id: userId });
    },
    redirectPath: "/dashboard/users",
    resourceName: "user",
  });
}
