import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { UserForm } from "./UserForm";
import type { UserFormSchema } from "../types";

interface UserFormCardProps {
  title: string;
  defaultValues?: Partial<UserFormSchema>;
  onSubmit: (data: UserFormSchema) => void;
  isSubmitting?: boolean;
  submitButtonText?: string;
  showPasswordField?: boolean;
}

export function UserFormCard({
  title,
  defaultValues,
  onSubmit,
  isSubmitting,
  submitButtonText,
  showPasswordField = true,
}: UserFormCardProps) {
  return (
    <div className="container w-full mx-auto p-5">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">{title}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <UserForm
            defaultValues={defaultValues}
            onSubmit={onSubmit}
            isSubmitting={isSubmitting}
            submitButtonText={submitButtonText}
            showPasswordField={showPasswordField}
          />
        </CardContent>
      </Card>
    </div>
  );
}
