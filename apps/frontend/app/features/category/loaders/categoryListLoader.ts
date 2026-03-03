import { createListLoader } from "~/lib/loaders";
import { categoryApi } from "../api";

export async function categoryListLoader({ request }: { request: Request }) {
  return createListLoader({
    request,
    getAllFn: categoryApi.getAll.bind(categoryApi),
    searchField: "title",
  });
}
