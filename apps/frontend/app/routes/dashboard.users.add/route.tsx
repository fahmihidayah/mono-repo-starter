import { createUserAction } from "~/features/users/actions";
import { UserFormCard } from "~/features/users/components";
import { useUserForm } from "~/features/users/hooks";

export async function action({ request }: { request: Request }) {
  return createUserAction({ request });
}

export default function AddUserPage() {
  const { handleSubmit, isSubmitting } = useUserForm(true);

  return (
    <UserFormCard
      title="Add New User"
      defaultValues={{ name: "", email: "", password: "" }}
      onSubmit={handleSubmit}
      isSubmitting={isSubmitting}
      submitButtonText="Submit"
      showPasswordField={true}
    />
  );
}
