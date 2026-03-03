import { redirect } from "react-router";
import { postApi } from "../api";
import { postFormSchema } from "../types";
import type { ActionData } from "~/types";
import { createResourceAction } from "~/lib/actions";

export async function createPostAction({ request }: { request: Request }) {
  return createResourceAction({
    request,
    schema: postFormSchema,
    createFn: async (data, req) => {
      return postApi.create({ request: req, data });
    },
    redirectPath: "/dashboard/posts",
    resourceName: "post",
    arrayFields: ["category_ids"],
  });
}

// Previous implementation without using createResourceAction helper
//   }
// }
// export async function createPostAction({ request }: { request: Request }) {
//   const formData = await request.formData();
//   const categoryIds = formData.getAll("categoryIds");
//   const data = {
//     title: formData.get("title"),
//     content: formData.get("content"),
//     category_ids: categoryIds,
//   };

//   const validation = postFormSchema.safeParse(data);

//   if (!validation.success) {
//     const fieldErrors = validation.error.flatten((issue) => issue.message).fieldErrors;
//     return {
//       success: false,
//       errors: {
//         ...fieldErrors,
//       },
//     } as ActionData;
//   }

//   try {
//     await postApi.create({
//       request,
//       data: {
//         title: validation.data.title,
//         content: validation.data.content,
//         category_ids: validation.data.category_ids,
//       },
//     });
//     return redirect("/dashboard/posts?success=created");
//   } catch (error) {
//     console.error("Error creating post:", error);
//     return {
//       success: false,
//       errors: {
//         general: "Failed to create post. Please try again.",
//       },
//     } as ActionData;
//   }
// }
