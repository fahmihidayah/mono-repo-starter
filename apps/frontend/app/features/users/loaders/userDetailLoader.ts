import { userApi } from "../api/users";

export async function userDetailLoader({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const user = await userApi.getById({
    request,
    id,
  });
  return { user };
}
