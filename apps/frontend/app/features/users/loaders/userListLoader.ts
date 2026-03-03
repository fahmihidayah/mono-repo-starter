import { createListLoader } from "~/lib/loaders";
import { userApi } from "../api/users";

export async function userListLoader({ request }: { request: Request }) {
  return createListLoader({
    request,
    getAllFn: userApi.getAll.bind(userApi),
    searchField: "email",
  });
}
