import { updateResourceAction } from "~/lib/actions";
import { categoryApi } from "../api";
import { categoryFormSchema } from "../types";

export async function updateCategoryAction({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  return updateResourceAction({
    request,
    id,
    schema: categoryFormSchema,
    updateFn: async (categoryId, data, req) => {
      return categoryApi.update({ request: req, data, id: categoryId });
    },
    redirectPath: "/dashboard/categories",
    resourceName: "category",
  });
}
