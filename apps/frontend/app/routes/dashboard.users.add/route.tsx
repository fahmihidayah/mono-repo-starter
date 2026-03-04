import { createUserAction } from "~/features/users/actions";
import { UserFormCard } from "~/features/users/components";
import { useUserForm } from "~/features/users/hooks";
import { useResourceForm } from "~/lib/hooks";

export async function action({ request }: { request: Request }) {
  return createUserAction({ request });
}

export default function AddUserPage() {
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <UserFormCard
      title="Add New User"
      defaultValues={{ name: "", email: "", password: "" }}
      onSubmit={submitFormSchema}
      isSubmitting={isSubmitting}
      submitButtonText="Submit"
      showPasswordField={true}
    />
  );
}
