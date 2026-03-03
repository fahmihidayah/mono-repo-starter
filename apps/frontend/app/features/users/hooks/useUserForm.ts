import { useResourceForm } from "~/lib/hooks";
import type { UserFormSchema } from "../types";

export function useUserForm(includePassword: boolean = true) {
  const { handleSubmit: baseHandleSubmit, isSubmitting, actionData } = useResourceForm();

  const handleSubmit = (data: UserFormSchema) => {
    const formData = new FormData();
    formData.append("name", data.name);
    formData.append("email", data.email);
    if (includePassword && data.password) {
      formData.append("password", data.password);
    }
    baseHandleSubmit(formData);
  };

  return {
    handleSubmit,
    isSubmitting,
    actionData,
  };
}
