interface DeleteResourceActionConfig {
  request: Request;
  deleteFn: (id: string, request: Request) => Promise<unknown>;
  deleteManyFn: (ids: string[], request: Request) => Promise<unknown>;
  resourceName: string;
}

export async function deleteResourceAction({
  request,
  deleteFn,
  deleteManyFn,
  resourceName,
}: DeleteResourceActionConfig) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const id = formData.get("id")?.toString();
  const ids = formData.getAll("ids[]");

  try {
    if (intent === "delete" && id) {
      await deleteFn(id, request);
      return {
        success: true,
        message: `${resourceName} deleted successfully`,
      };
    }

    if (intent === "bulkDelete" && ids.length > 0) {
      await deleteManyFn(
        ids.map((id) => id.toString()),
        request
      );
      const count = ids.length;
      const pluralName = count !== 1 ? `${resourceName}s` : resourceName;
      return {
        success: true,
        message: `${count} ${pluralName.toLowerCase()} deleted successfully`,
      };
    }

    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}
