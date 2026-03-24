import { createListLoader } from "~/lib/loaders";
import { roleApi } from "../api";

export async function roleListLoader({ request }: { request: Request }) {
  return createListLoader({
    request,
    getAllFn: roleApi.getAll.bind(roleApi),
    searchField: "name",
  });
}
