import { redirect } from "react-router";
import { userApi } from "../api/users";
import { userFormSchema } from "../types";
import type { ActionData } from "~/types";

export async function updateUserAction({
  request,
  id,
}: {
  request: Request;
  id: string;
}) {
  const formData = await request.formData();
  const data = Object.fromEntries(formData);
  const validation = userFormSchema.safeParse(data);

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
    await userApi.update({
      request,
      data: data,
      id,
    });
    return redirect("/dashboard/users");
  } catch (error) {
    return {
      success: false,
      errors: {
        general: "Failed to update user. Please try again.",
      },
    } as ActionData;
  }
}
