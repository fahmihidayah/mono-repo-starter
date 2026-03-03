import { redirect } from "react-router";
import { userApi } from "../api/users";
import { userFormSchema } from "../types";
import type { ActionData } from "~/types";

export async function createUserAction({ request }: { request: Request }) {
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
    await userApi.create({
      request,
      data: data,
    });
    return redirect("/dashboard/users");
  } catch (error) {
    return {
      success: false,
      errors: {
        general: "Failed to create user. Please try again.",
      },
    } as ActionData;
  }
}
