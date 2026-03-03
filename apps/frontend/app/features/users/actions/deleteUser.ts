import { userApi } from "../api/users";

export async function deleteUserAction({ request }: { request: Request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const id = formData.get("id")?.toString();
  const ids = formData.getAll("ids[]");

  try {
    if (intent === "delete" && id) {
      await userApi.deleteById({
        request,
        id,
      });
      return { success: true, message: "User deleted successfully" };
    }

    if (intent === "bulkDelete" && ids.length > 0) {
      await userApi.deleteMany({
        request,
        ids: ids.map((id) => id.toString()),
      });
      return {
        success: true,
        message: `${ids.length} user${ids.length !== 1 ? "s" : ""} deleted successfully`,
      };
    }

    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}
