import { createResourceAction } from "~/lib/actions";
import { categoryApi } from "../api";
import { categoryFormSchema } from "../types";

export async function createCategoryAction({ request }: { request: Request }) {
  return createResourceAction({
    request,
    schema: categoryFormSchema,
    createFn: async (data, req) => {
      return categoryApi.create({ request: req, data });
    },
    redirectPath: "/dashboard/categories",
    resourceName: "category",
  });
}
