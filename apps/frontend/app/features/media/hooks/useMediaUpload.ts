import { useSubmit, useNavigation, useActionData } from "react-router";
import { useEffect } from "react";
import { toast } from "sonner";

type ActionData = {
  success: boolean;
  errors?: Record<string, string[]>;
};

export function useMediaUpload() {
  const submit = useSubmit();
  const navigation = useNavigation();
  const actionData = useActionData<ActionData>();
  const isSubmitting = navigation.state === "submitting";

  useEffect(() => {
    if (actionData?.errors?._form) {
      toast.error(actionData.errors._form[0]);
    }
  }, [actionData]);

  const handleUpload = (files: File[], alt: string) => {
    const formData = new FormData();

    files.forEach((file) => {
      formData.append("files", file);
    });

    if (alt) {
      formData.append("alt", alt);
    }

    submit(formData, { method: "post", encType: "multipart/form-data" });
  };

  return {
    handleUpload,
    isSubmitting,
    actionData,
  };
}
