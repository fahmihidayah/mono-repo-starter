import { useSubmit, useActionData, useNavigation } from "react-router";
import { useEffect } from "react";
import { toast } from "sonner";
import type { CategoryFormSchema } from "../types";
import type { ActionData } from "~/types";

export function useCategoryForm() {
  const submit = useSubmit();
  const actionData = useActionData<ActionData>();
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

  useEffect(() => {
    if (actionData && !actionData.success) {
      if (actionData.errors?.general) {
        toast.error(actionData.errors.general);
      }
    }
  }, [actionData]);

  const handleSubmit = (data: CategoryFormSchema) => {
    const formData = new FormData();
    formData.append("title", data.title);
    submit(formData, { method: "post" });
  };

  return {
    handleSubmit,
    isSubmitting,
    actionData,
  };
}
