import { deleteResourceAction } from "~/lib/actions";
import { mediaApi } from "../api";

export async function deleteMediaAction({ request }: { request: Request }) {
  return deleteResourceAction({
    request,
    deleteFn: async (id, req) => {
      return mediaApi.deleteById({ request: req, id });
    },
    deleteManyFn: async (ids, req) => {
      return mediaApi.deleteMany({ request: req, ids });
    },
    resourceName: "media",
  });
}
