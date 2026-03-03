import { redirect } from "react-router";
import type { z } from "zod";
import type { ActionData } from "~/types";
import { formDataToObject } from "../utils/formData";

interface UpdateResourceActionConfig<T extends z.ZodTypeAny> {
  request: Request;
  id: string;
  schema: T;
  updateFn: (id: string, data: z.infer<T>, request: Request) => Promise<unknown>;
  redirectPath: string;
  resourceName: string;
  arrayFields?: string[];
}

export async function updateResourceAction<T extends z.ZodTypeAny>({
  request,
  id,
  schema,
  updateFn,
  redirectPath,
  resourceName,
  arrayFields = [],
}: UpdateResourceActionConfig<T>) {
  const formData = await request.formData();
  const data = formDataToObject(formData, arrayFields);
  const validation = schema.safeParse(data);

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
    await updateFn(id, validation.data, request);
    return redirect(redirectPath);
  } catch (error) {
    console.error(`Error updating ${resourceName}:`, error);
    return {
      success: false,
      errors: {
        general: `Failed to update ${resourceName}. Please try again.`,
      },
    } as ActionData;
  }
}
