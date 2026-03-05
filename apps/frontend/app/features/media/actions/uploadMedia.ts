import { redirect } from "react-router";
import { getToken } from "~/lib/utils.server";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export async function uploadMediaAction({ request }: { request: Request }) {
  const formData = await request.formData();
  const token = await getToken({ request });

  // Get all files from the form
  const files = formData.getAll("files") as File[];
  const alt = formData.get("alt") as string;

  if (files.length === 0) {
    return {
      success: false,
      errors: { _form: ["No files selected"] },
    };
  }

  try {
    let successCount = 0;
    let failCount = 0;

    // Upload each file separately (backend accepts one file at a time)
    for (const file of files) {
      const fileFormData = new FormData();
      fileFormData.append("media", file);
      if (alt) {
        fileFormData.append("alt", alt);
      }

      try {
        const response = await fetch(`${API_BASE_URL}/api/media`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: fileFormData,
        });

        if (response.ok) {
          successCount++;
        } else {
          failCount++;
          const errorData = await response.json();
          console.error(`Failed to upload ${file.name}:`, errorData.message);
        }
      } catch (error) {
        failCount++;
        console.error(`Failed to upload ${file.name}:`, error);
      }
    }

    if (successCount > 0) {
      return redirect("/dashboard/media?success=uploaded");
    }

    return {
      success: false,
      errors: { _form: [`Failed to upload ${failCount} file${failCount > 1 ? "s" : ""}`] },
    };
  } catch (error) {
    console.error("Upload error:", error);
    return {
      success: false,
      errors: { _form: ["Failed to upload media files"] },
    };
  }
}
