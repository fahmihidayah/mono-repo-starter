import { createDetailLoader } from "~/lib/loaders";
import { categoryApi } from "../api";

export async function categoryDetailLoader({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const category = await createDetailLoader({
    request,
    id,
    getByIdFn: async (categoryId, req) => {
      return categoryApi.getById({ request: req, id: categoryId });
    },
  });

  return { category };
}
