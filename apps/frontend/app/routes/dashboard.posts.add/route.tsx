import { useLoaderData } from "react-router";
import { postFormLoader } from "~/features/posts/loaders";
import { createPostAction } from "~/features/posts/actions";
import { PostForm } from "~/features/posts/components";
import { useResourceForm } from "~/lib/hooks";

export async function loader({ request }: { request: Request }) {
  return postFormLoader({ request });
}

export async function action({ request }: { request: Request }) {
  return createPostAction({ request });
}

export default function DashboardAddPostPage() {
  const { categories } = useLoaderData<typeof loader>();
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <div className="container w-full mx-auto">
      <div className="flex-1">
        <PostForm
          defaultValues={{ title: "", content: "", category_ids: [] }}
          categories={categories}
          onSubmit={submitFormSchema}
          isSubmitting={isSubmitting}
          submitButtonText="Save"
        />
      </div>
    </div>
  );
}
