import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { RoleForm } from "./RoleForm";
import type { RoleFormSchema } from "../types";

interface RoleFormCardProps {
  title: string;
  defaultValues?: Partial<RoleFormSchema>;
  onSubmit: (data: RoleFormSchema) => void;
  isSubmitting?: boolean;
  submitButtonText?: string;
}

export function RoleFormCard({
  title,
  defaultValues,
  onSubmit,
  isSubmitting,
  submitButtonText,
}: RoleFormCardProps) {
  return (
    <div className="container w-full mx-auto p-5">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">{title}</CardTitle>
        </CardHeader>
        <CardContent>
          <RoleForm
            defaultValues={defaultValues}
            onSubmit={onSubmit}
            isSubmitting={isSubmitting}
            submitButtonText={submitButtonText}
          />
        </CardContent>
      </Card>
    </div>
  );
}
