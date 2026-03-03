import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { CategoryForm } from "./CategoryForm";
import type { CategoryFormSchema } from "../types";

interface CategoryFormCardProps {
  title: string;
  defaultValues?: Partial<CategoryFormSchema>;
  onSubmit: (data: CategoryFormSchema) => void;
  isSubmitting?: boolean;
  submitButtonText?: string;
}

export function CategoryFormCard({
  title,
  defaultValues,
  onSubmit,
  isSubmitting,
  submitButtonText,
}: CategoryFormCardProps) {
  return (
    <div className="container w-full mx-auto p-5">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">{title}</CardTitle>
        </CardHeader>
        <CardContent>
          <CategoryForm
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
