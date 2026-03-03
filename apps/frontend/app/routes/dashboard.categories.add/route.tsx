import { createCategoryAction } from "~/features/category/actions";
import { CategoryFormCard } from "~/features/category/components";
import { useCategoryForm } from "~/features/category/hooks";

export async function action({ request }: { request: Request }) {
  return createCategoryAction({ request });
}

export default function DashboardAddCategoryPage() {
  const { handleSubmit, isSubmitting } = useCategoryForm();

  return (
    <CategoryFormCard
      title="Add New Category"
      defaultValues={{ title: "" }}
      onSubmit={handleSubmit}
      isSubmitting={isSubmitting}
      submitButtonText="Submit"
    />
  );
}
