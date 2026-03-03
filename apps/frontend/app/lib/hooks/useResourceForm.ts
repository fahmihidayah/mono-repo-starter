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

  const handleSubmit = (formData: FormData) => {
    if (config?.onSubmit) {
      config.onSubmit(formData);
    }
    submit(formData, { method: "post" });
  };

  return {
    handleSubmit,
    isSubmitting,
    actionData,
  };
}
