import { useResourceForm } from "~/lib/hooks";
import { appendToFormData } from "~/lib/utils/formData";
import type { PostFormSchema } from "../types";

export function usePostForm() {
  const { handleSubmit: baseHandleSubmit, isSubmitting, actionData } = useResourceForm();

  const handleSubmit = (data: PostFormSchema) => {
    const formData = new FormData();
    formData.append("title", data.title);
    formData.append("content", data.content);
    appendToFormData(formData, data, "category_ids");
    baseHandleSubmit(formData);
  };

  return {
    handleSubmit,
    isSubmitting,
    actionData,
  };
}
