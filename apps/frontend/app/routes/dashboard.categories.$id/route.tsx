import { useLoaderData } from "react-router";
import type { Route } from "./+types/route";
import { categoryDetailLoader } from "~/features/category/loaders";
import { updateCategoryAction } from "~/features/category/actions";
import { CategoryFormCard } from "~/features/category/components";
import { useCategoryForm } from "~/features/category/hooks";

export async function loader({ request, params }: Route.LoaderArgs) {
  return categoryDetailLoader({ request, id: params.id });
}

export async function action({ request, params }: Route.ActionArgs) {
  return updateCategoryAction({ request, id: params.id });
}

export default function EditCategoryPage() {
  const { category } = useLoaderData<typeof loader>();
  const { handleSubmit, isSubmitting } = useCategoryForm();

  return (
    <CategoryFormCard
      title="Edit Category"
      defaultValues={{
        title: category.data?.title || "",
      }}
      onSubmit={handleSubmit}
      isSubmitting={isSubmitting}
      submitButtonText="Update"
    />
  );
}
