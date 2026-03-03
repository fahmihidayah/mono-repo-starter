import { userApi } from "../api/users";

export async function userListLoader({ request }: { request: Request }) {
  const url = new URL(request.url);
  const page = Number.parseInt(url.searchParams.get("page") || "1");
  const pageSize = Number.parseInt(url.searchParams.get("limit") || "10");
  const search = url.searchParams.get("search") || "";

  const result = await userApi.getAll({
    request,
    page,
    pageSize,
    search: { email: search },
  });

  return result;
}
