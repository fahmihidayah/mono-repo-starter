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
import { roleFormSchema, type RoleFormSchema } from "../types";

interface RoleFormProps {
  defaultValues?: Partial<RoleFormSchema>;
  onSubmit: (data: RoleFormSchema) => void;
  isSubmitting?: boolean;
  submitButtonText?: string;
}

export function RoleForm({
  defaultValues = { name: "" },
  onSubmit,
  isSubmitting = false,
  submitButtonText = "Submit",
}: RoleFormProps) {
  const form = useForm<RoleFormSchema>({
    resolver: zodResolver(roleFormSchema),
    defaultValues,
  });

  return (
    <Form {...form}>
      <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
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
