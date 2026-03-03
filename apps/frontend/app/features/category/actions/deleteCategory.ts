import { categoryApi } from "../api";

export async function deleteCategoryAction({ request }: { request: Request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const id = formData.get("id")?.toString();
  const ids = formData.getAll("ids[]");

  try {
    if (intent === "delete" && id) {
      await categoryApi.deleteById({
        request,
        id,
      });
      return { success: true, message: "Category deleted successfully" };
    }

    if (intent === "bulkDelete" && ids.length > 0) {
      await categoryApi.deleteMany({
        request,
        ids: ids.map((id) => id.toString()),
      });
      return {
        success: true,
        message: `${ids.length} categor${ids.length !== 1 ? "ies" : "y"} deleted successfully`,
      };
    }

    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}
