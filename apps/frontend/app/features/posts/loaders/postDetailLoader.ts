import { postApi } from "../api";
import { categoryApi } from "~/features/category/api";

export async function postDetailLoader({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const [post, categoriesResponse] = await Promise.all([
    postApi.getById({ request, id }),
    categoryApi.getAll({ request, page: 1, pageSize: 100 }),
  ]);

  return {
    post: post.data,
    categories: categoriesResponse.data || [],
  };
}
