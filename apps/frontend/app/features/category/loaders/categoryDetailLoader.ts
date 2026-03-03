import { categoryApi } from "../api";

export async function categoryDetailLoader({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const category = await categoryApi.getById({
    request,
    id,
  });
  return { category };
}
