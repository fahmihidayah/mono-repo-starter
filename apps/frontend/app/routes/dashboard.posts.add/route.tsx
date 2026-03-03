import { useLoaderData } from "react-router";
import { postFormLoader } from "~/features/posts/loaders";
import { createPostAction } from "~/features/posts/actions";
import { PostForm } from "~/features/posts/components";
import { usePostForm } from "~/features/posts/hooks";

export async function loader({ request }: { request: Request }) {
  return postFormLoader({ request });
}

export async function action({ request }: { request: Request }) {
  return createPostAction({ request });
}

export default function DashboardAddPostPage() {
  const { categories } = useLoaderData<typeof loader>();
  const { handleSubmit, isSubmitting } = usePostForm();

  return (
    <div className="container w-full mx-auto">
      <div className="flex-1">
        <PostForm
          defaultValues={{ title: "", content: "", category_ids: [] }}
          categories={categories}
          onSubmit={handleSubmit}
          isSubmitting={isSubmitting}
          submitButtonText="Save"
        />
      </div>
    </div>
  );
}
