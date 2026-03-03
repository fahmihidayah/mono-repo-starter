import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Form as ReactRouterForm } from "react-router";
import { Button } from "~/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { categoryFormSchema, type CategoryFormSchema } from "../types";

interface CategoryFormProps {
  defaultValues?: Partial<CategoryFormSchema>;
  onSubmit: (data: CategoryFormSchema) => void;
  isSubmitting?: boolean;
  submitButtonText?: string;
}

export function CategoryForm({
  defaultValues = { title: "" },
  onSubmit,
  isSubmitting = false,
  submitButtonText = "Submit",
}: CategoryFormProps) {
  const form = useForm<CategoryFormSchema>({
    resolver: zodResolver(categoryFormSchema),
    defaultValues,
  });

  return (
    <Form {...form}>
      <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Title</FormLabel>
              <FormControl>
                <Input disabled={field.disabled || isSubmitting} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Submitting..." : submitButtonText}
        </Button>
      </ReactRouterForm>
    </Form>
  );
}
