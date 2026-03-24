import { createRoleAction } from "~/features/roles/actions";
import { RoleFormCard } from "~/features/roles/components";
import { useResourceForm } from "~/lib/hooks";

export async function action({ request }: { request: Request }) {
  return createRoleAction({ request });
}

export default function DashboardAddRolePage() {
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <RoleFormCard
      title="Add New Role"
      defaultValues={{ name: "" }}
      onSubmit={submitFormSchema}
      isSubmitting={isSubmitting}
      submitButtonText="Submit"
    />
  );
}
