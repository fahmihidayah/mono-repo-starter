import { createResourceAction } from "~/lib/actions";
import { userApi } from "../api/users";
import { userFormSchema } from "../types";

export async function createUserAction({ request }: { request: Request }) {
  return createResourceAction({
    request,
    schema: userFormSchema,
    createFn: async (data, req) => {
      return userApi.create({ request: req, data });
    },
    redirectPath: "/dashboard/users",
    resourceName: "user",
  });
}
