import { useSubmit, useActionData, useNavigation } from "react-router";
import { useEffect } from "react";
import { toast } from "sonner";
import type { UserFormSchema } from "../types";
import type { ActionData } from "~/types";

export function useUserForm(includePassword: boolean = true) {
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

  const handleSubmit = (data: UserFormSchema) => {
    const formData = new FormData();
    formData.append("name", data.name);
    formData.append("email", data.email);
    if (includePassword && data.password) {
      formData.append("password", data.password);
    }
    submit(formData, { method: "post" });
  };

  return {
    handleSubmit,
    isSubmitting,
    actionData,
  };
}
