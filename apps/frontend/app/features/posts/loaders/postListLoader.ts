import { createListLoader } from "~/lib/loaders";
import { postApi } from "../api";

export async function postListLoader({ request }: { request: Request }) {
  return createListLoader({
    request,
    getAllFn: postApi.getAll.bind(postApi),
    searchField: "title",
  });
}
