import { categoryApi } from "~/features/category/api";

export async function postFormLoader({ request }: { request: Request }) {
  const categories = await categoryApi.getAll({
    request,
    page: 1,
    pageSize: 100,
  });

  return { categories: categories.data || [] };
}
