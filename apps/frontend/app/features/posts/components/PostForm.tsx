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
import { MultiSelect } from "~/components/ui/multi-select";
import LexicalEditorField from "~/components/editor/LexicalEditorField";
import { postFormSchema, type PostFormSchema } from "../types";
import type { Category } from "~/features/category/types";

interface PostFormProps {
  defaultValues?: Partial<PostFormSchema>;
  categories: Category[];
  onSubmit: (data: PostFormSchema) => void;
  isSubmitting?: boolean;
  submitButtonText?: string;
}

export function PostForm({
  defaultValues = { title: "", content: "", category_ids: [] },
  categories,
  onSubmit,
  isSubmitting = false,
  submitButtonText = "Submit",
}: PostFormProps) {
  const form = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
    defaultValues,
  });

  return (
    <Form {...form}>
      <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 p-5">
        <div className="flex flex-row justify-between items-center">
          <h3 className="text-2xl font-bold">
            {submitButtonText === "Update" ? "Edit Post" : "Add New Post"}
          </h3>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Submitting..." : submitButtonText}
          </Button>
        </div>
        <hr className="border-b" />

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

        <FormField
          control={form.control}
          name="content"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Content</FormLabel>
              <FormControl>
                <LexicalEditorField
                  value={field.value}
                  onChange={field.onChange}
                  placeholder="Write your post content..."
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="category_ids"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Categories</FormLabel>
              <FormControl>
                <MultiSelect
                  options={categories.map((cat) => ({
                    label: cat.title || "",
                    value: cat.id,
                  }))}
                  value={field.value}
                  onChange={field.onChange}
                  placeholder="Select categories..."
                  disabled={field.disabled || isSubmitting}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </ReactRouterForm>
    </Form>
  );
}
