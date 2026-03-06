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
import { ImageField } from "~/components/ui/image-field";
import { mediaFormSchema, type MediaFormSchema } from "../types";

interface MediaFormProps {
  defaultValues?: Partial<MediaFormSchema>;
  onSubmit: (data: MediaFormSchema) => void;
  isSubmitting?: boolean;
  submitButtonText?: string;
}

export function MediaForm({
  defaultValues = { alt: "" },
  onSubmit,
  isSubmitting = false,
  submitButtonText = "Submit",
}: MediaFormProps) {
  const form = useForm<MediaFormSchema>({
    resolver: zodResolver(mediaFormSchema),
    defaultValues,
  });

  return (
    <Form {...form}>
      <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="image"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Image</FormLabel>
              <FormControl>
                <ImageField
                  value={field.value || null}
                  onChange={field.onChange}
                  disabled={field.disabled || isSubmitting}
                  accept="image/*"
                  maxSize={5 * 1024 * 1024}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="alt"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Alt Text (Optional)</FormLabel>
              <FormControl>
                <Input
                  disabled={field.disabled || isSubmitting}
                  placeholder="Describe the image for accessibility"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Uploading..." : submitButtonText}
        </Button>
      </ReactRouterForm>
    </Form>
  );
}
