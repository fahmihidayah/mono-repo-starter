import { useResourceForm } from "~/lib/hooks";
import type { CategoryFormSchema } from "../types";

export function useCategoryForm() {
  const { handleSubmit: baseHandleSubmit, isSubmitting, actionData } = useResourceForm();

  const handleSubmit = (data: CategoryFormSchema) => {
    const formData = new FormData();
    formData.append("title", data.title);
    baseHandleSubmit(formData);
  };

  return {
    handleSubmit,
    isSubmitting,
    actionData,
  };
}
