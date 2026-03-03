import { redirect } from "react-router";
import type { z } from "zod";
import type { ActionData } from "~/types";
import { formDataToObject } from "../utils/formData";

interface CreateResourceActionConfig<T extends z.ZodTypeAny> {
  request: Request;
  schema: T;
  createFn: (data: z.infer<T>, request: Request) => Promise<unknown>;
  redirectPath: string;
  resourceName: string;
  arrayFields?: string[];
}

export async function createResourceAction<T extends z.ZodTypeAny>({
  request,
  schema,
  createFn,
  redirectPath,
  resourceName,
  arrayFields = [],
}: CreateResourceActionConfig<T>) {
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
    await createFn(validation.data, request);
    return redirect(`${redirectPath}?success=created`);
  } catch (error) {
    console.error(`Error creating ${resourceName}:`, error);
    return {
      success: false,
      errors: {
        general: `Failed to create ${resourceName}. Please try again.`,
      },
    } as ActionData;
  }
}
