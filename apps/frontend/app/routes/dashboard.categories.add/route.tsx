import { createCategoryAction } from "~/features/category/actions";
import { CategoryFormCard } from "~/features/category/components";
import { useResourceForm } from "~/lib/hooks";

export async function action({ request }: { request: Request }) {
  return createCategoryAction({ request });
}

export default function DashboardAddCategoryPage() {
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <CategoryFormCard
      title="Add New Category"
      defaultValues={{ title: "" }}
      onSubmit={submitFormSchema}
      isSubmitting={isSubmitting}
      submitButtonText="Submit"
    />
  );
}
