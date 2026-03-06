import { redirect } from "react-router";
import type { ActionData } from "~/types";
import { formDataToObject } from "~/lib/utils/formData";
import { mediaApi } from "../api";
import { mediaFormSchema } from "../types";

export async function createMediaAction({ request }: { request: Request }) {
  const formData = await request.formData();
  const data = formDataToObject(formData);

  // Log received data for debugging
  console.log("Received form data:", {
    keys: Array.from(formData.keys()),
    image: data.image,
    imageType: data.image instanceof File ? "File" : typeof data.image,
    alt: data.alt,
  });

  // Validate the form data
  const validation = mediaFormSchema.safeParse(data);

  if (!validation.success) {
    const fieldErrors = validation.error.flatten((issue) => issue.message).fieldErrors;
    console.error("Media validation failed:", fieldErrors);

    // Convert field errors from string[] to string
    const errors: Record<string, string> = {};
    Object.entries(fieldErrors).forEach(([key, value]) => {
      if (value && value.length > 0) {
        errors[key] = value[0]; // Take the first error message
      }
    });
    errors.general = "Please check the form and try again.";

    return {
      success: false,
      errors,
    } as ActionData;
  }

  const { image, alt } = validation.data;

  try {
    // Prepare FormData for API (backend expects "media" as field name)
    const apiFormData = new FormData();
    apiFormData.append("media", image);
    if (alt) {
      apiFormData.append("alt", alt);
    }

    console.log("Uploading media:", { fileName: image.name, size: image.size, alt });

    const result = await mediaApi.create({ request, data: apiFormData });

    if (result.code === 200 || result.code === 201) {
      console.log("Media uploaded successfully:", result.data);
      return redirect("/dashboard/media?success=uploaded");
    }

    // API returned an error response
    console.error("Media upload failed:", { code: result.code, message: result.message });
    return {
      success: false,
      errors: {
        general: result.message || "Failed to upload image. Please try again.",
      },
    } as ActionData;
  } catch (error) {
    console.error("Error uploading media:", error);
    return {
      success: false,
      errors: {
        general:
          error instanceof Error ? error.message : "Failed to upload image. Please try again.",
      },
    } as ActionData;
  }
}
