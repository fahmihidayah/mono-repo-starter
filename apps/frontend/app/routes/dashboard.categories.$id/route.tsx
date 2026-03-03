import { useLoaderData } from "react-router";
import type { Route } from "./+types/route";
import { categoryDetailLoader } from "~/features/category/loaders";
import { updateCategoryAction } from "~/features/category/actions";
import { CategoryFormCard } from "~/features/category/components";
import { useResourceForm } from "~/lib/hooks";

export async function loader({ request, params }: Route.LoaderArgs) {
  return categoryDetailLoader({ request, id: params.id });
}

export async function action({ request, params }: Route.ActionArgs) {
  return updateCategoryAction({ request, id: params.id });
}

export default function EditCategoryPage() {
  const { category } = useLoaderData<typeof loader>();
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <CategoryFormCard
      title="Edit Category"
      defaultValues={{
        title: category.data?.title || "",
      }}
      onSubmit={submitFormSchema}
      isSubmitting={isSubmitting}
      submitButtonText="Update"
    />
  );
}
