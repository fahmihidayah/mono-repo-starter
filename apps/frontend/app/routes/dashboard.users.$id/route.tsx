import { useLoaderData } from "react-router";
import type { Route } from "./+types/route";
import { userDetailLoader } from "~/features/users/loaders";
import { updateUserAction } from "~/features/users/actions";
import { UserFormCard } from "~/features/users/components";
import { useUserForm } from "~/features/users/hooks";

export async function loader({ request, params }: Route.LoaderArgs) {
  return userDetailLoader({ request, id: params.id });
}

export async function action({ request, params }: Route.ActionArgs) {
  return updateUserAction({ request, id: params.id });
}

export default function EditUserPage() {
  const { user } = useLoaderData<typeof loader>();
  const { handleSubmit, isSubmitting } = useUserForm(false);

  return (
    <UserFormCard
      title="Edit User"
      defaultValues={{
        name: user.data?.name || "",
        email: user.data?.email || "",
        password: "",
      }}
      onSubmit={handleSubmit}
      isSubmitting={isSubmitting}
      submitButtonText="Submit"
      showPasswordField={false}
    />
  );
}
