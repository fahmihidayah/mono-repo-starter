import { useSubmit, useActionData, useNavigation } from "react-router";
import { useEffect } from "react";
import { toast } from "sonner";
import type { ActionData } from "~/types";

interface UseResourceFormConfig {
  onSubmit?: (formData: FormData) => void;
}

export function useResourceForm(config?: UseResourceFormConfig) {
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

  const submitFormSchema = (form: any) => {
    const formData = new FormData();
    Object.entries(form).forEach(([key, value]) => {
      if (Array.isArray(value)) {
        formData.delete(key);
        value.forEach((item) => formData.append(key, item));
      } else if (value instanceof File) {
        formData.append(key, value);
      } else if (typeof value === "string") {
        formData.append(key, value);
      } else if (typeof value === "number" || typeof value === "boolean") {
        formData.append(key, value.toString());
      } else {
        formData.append(key, `${value}`);
      }
    });
    formData;
    handleSubmit(formData);
  };

  const handleSubmit = (formData: FormData) => {
    if (config?.onSubmit) {
      config.onSubmit(formData);
    }
    submit(formData, { method: "post" });
  };

  return {
    submitFormSchema,
    handleSubmit,
    isSubmitting,
    actionData,
  };
}
