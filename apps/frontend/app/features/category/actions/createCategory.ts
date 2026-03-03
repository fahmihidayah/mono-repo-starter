import { redirect } from "react-router";
import { categoryApi } from "../api";
import { categoryFormSchema } from "../types";
import type { ActionData } from "~/types";

export async function createCategoryAction({ request }: { request: Request }) {
  const formData = await request.formData();
  const data = Object.fromEntries(formData);
  const validation = categoryFormSchema.safeParse(data);

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
    await categoryApi.create({
      request,
      data: data,
    });
    return redirect("/dashboard/categories?success=created");
  } catch (error) {
    return {
      success: false,
      errors: {
        general: "Failed to create category. Please try again.",
      },
    } as ActionData;
  }
}
