import { redirect } from "react-router";
import { postApi } from "../api";
import { postFormSchema } from "../types";
import type { ActionData } from "~/types";

export async function updatePostAction({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const formData = await request.formData();
  const categoryIds = formData.getAll("category_ids");
  const data = {
    title: formData.get("title"),
    content: formData.get("content"),
    category_ids: categoryIds,
  };

  const validation = postFormSchema.safeParse(data);

  if (!validation.success) {
    const fieldErrors = validation.error.flatten((issue) => issue.message).fieldErrors;
    return {
      success: false,
      errors: {
        ...fieldErrors,
      },
    } as ActionData;
  }

  try {
    await postApi.update({
      request,
      data: {
        title: validation.data.title,
        content: validation.data.content,
        category_ids: validation.data.category_ids,
      },
      id,
    });
    return redirect("/dashboard/posts");
  } catch (error) {
    return {
      success: false,
      errors: {
        general: "Failed to update post. Please try again.",
      },
    } as ActionData;
  }
}
