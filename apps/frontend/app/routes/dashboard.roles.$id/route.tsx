import { useLoaderData } from "react-router";
import type { Route } from "./+types/route";
import { roleDetailLoader } from "~/features/roles/loaders";
import { updateRoleAction } from "~/features/roles/actions";
import { RoleFormCard } from "~/features/roles/components";
import { useResourceForm } from "~/lib/hooks";

export async function loader({ request, params }: Route.LoaderArgs) {
  return roleDetailLoader({ request, id: params.id });
}

export async function action({ request, params }: Route.ActionArgs) {
  return updateRoleAction({ request, id: params.id });
}

export default function EditRolePage() {
  const { role } = useLoaderData<typeof loader>();
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <RoleFormCard
      title="Edit Role"
      defaultValues={{
        name: role.data?.name || "",
      }}
      onSubmit={submitFormSchema}
      isSubmitting={isSubmitting}
      submitButtonText="Update"
    />
  );
}
