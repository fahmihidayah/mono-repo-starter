import { deleteResourceAction } from "~/lib/actions";
import { postApi } from "../api";

export async function deletePostAction({ request }: { request: Request }) {
  return deleteResourceAction({
    request,
    deleteFn: async (id, req) => {
      return postApi.deleteById({ request: req, id });
    },
    deleteManyFn: async (ids, req) => {
      return postApi.deleteMany({ request: req, ids });
    },
    resourceName: "post",
  });
}

// export async function deletePostAction({ request }: { request: Request }) {
//   const formData = await request.formData();
//   const intent = formData.get("intent");
//   const id = formData.get("id")?.toString();
//   const ids = formData.getAll("ids[]");

//   try {
//     if (intent === "delete" && id) {
//       await postApi.deleteById({
//         request,
//         id,
//       });
//       return { success: true, message: "Post deleted successfully" };
//     }

//     if (intent === "bulkDelete" && ids.length > 0) {
//       await postApi.deleteMany({
//         request,
//         ids: ids.map((id) => id.toString()),
//       });
//       return {
//         success: true,
//         message: `${ids.length} post${ids.length !== 1 ? "s" : ""} deleted successfully`,
//       };
//     }

//     return { success: false, message: "Invalid action" };
//   } catch (error) {
//     console.error("Action error:", error);
//     return { success: false, message: "An error occurred" };
//   }
// }
