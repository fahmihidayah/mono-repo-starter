import { deleteResourceAction } from "~/lib/actions";
import { categoryApi } from "../api";

export async function deleteCategoryAction({ request }: { request: Request }) {
  return deleteResourceAction({
    request,
    deleteFn: async (id, req) => {
      return categoryApi.deleteById({ request: req, id });
    },
    deleteManyFn: async (ids, req) => {
      return categoryApi.deleteMany({ request: req, ids });
    },
    resourceName: "category",
  });
}
