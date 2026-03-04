import { useLoaderData } from "react-router";
import type { Route } from "./+types/route";
import { postDetailLoader } from "~/features/posts/loaders";
import { updatePostAction } from "~/features/posts/actions";
import { PostForm } from "~/features/posts/components";
import { useResourceForm } from "~/lib/hooks";

export async function loader({ request, params }: Route.LoaderArgs) {
  return postDetailLoader({ request, id: params.id });
}

export async function action({ request, params }: Route.ActionArgs) {
  return updatePostAction({ request, id: params.id });
}

export default function EditPostPage() {
  const { post, categories } = useLoaderData<typeof loader>();
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <div className="container w-full mx-auto p-5">
      <PostForm
        defaultValues={{
          title: post?.title || "",
          content: post?.content || "",
          category_ids: post?.categories?.map((cat) => cat.id) || [],
        }}
        categories={categories}
        onSubmit={submitFormSchema}
        isSubmitting={isSubmitting}
        submitButtonText="Update"
      />
    </div>
  );
}
