import { createListLoader } from "~/lib/loaders";
import { mediaApi } from "../api";

export async function mediaListLoader({ request }: { request: Request }) {
  return createListLoader({
    request,
    getAllFn: mediaApi.getAll.bind(mediaApi),
    searchField: "file_name",
  });
}
