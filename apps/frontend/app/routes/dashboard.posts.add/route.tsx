import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {
  Form as ReactRouterForm,
  redirect,
  useSubmit,
  useActionData,
  useNavigation,
  useLoaderData,
} from "react-router";
import { z } from "zod";
import { useEffect } from "react";
import { toast } from "sonner";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import { MultiSelect } from "~/components/ui/multi-select";
import LexicalEditorField from "~/components/editor/LexicalEditorField";
import { postApi } from "~/features/posts/api";
import { categoryApi } from "~/features/category/api";
import type { ActionData } from "~/types";
import { postFormSchema, type PostFormSchema } from "~/features/posts/types";

export async function loader({ request }: { request: Request }) {
  const categories = await categoryApi.getAll({
    request,
    page: 1,
    pageSize: 100,
  });

  return { categories: categories.data || [] };
}

export async function action({ request }: { request: Request }) {
  const formData = await request.formData();
  const categoryIds = formData.getAll("categoryIds");
  const data = {
    title: formData.get("title"),
    content: formData.get("content"),
    category_ids: categoryIds,
  };

  const validation = postFormSchema.safeParse(data);

  if (!validation.success) {
    const fieldErrors = validation.error.flatten((issue) => issue.message).fieldErrors;
    return {
      success: false,
      errors: {
        ...fieldErrors,
      },
    } as ActionData;
  }

  try {
    await postApi.create({
      request,
      data: {
        title: validation.data.title,
        content: validation.data.content,
        category_ids: validation.data.category_ids,
      },
    });
    return redirect("/dashboard/posts?success=created");
  } catch (error) {
    console.error("Error creating post:", error);
    return {
      success: false,
      errors: {
        general: "Failed to create post. Please try again.",
      },
    } as ActionData;
  }
}

export default function DashboardAddPostPage() {
  const { categories } = useLoaderData<typeof loader>();
  const form = useForm<PostFormSchema>({
    resolver: zodResolver(postFormSchema),
    defaultValues: {
      title: "",
      content: "",
      category_ids: [],
    },
  });

  const submit = useSubmit();
  const actionData = useActionData<ActionData>();
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

  useEffect(() => {
    if (actionData && !actionData.success) {
      if (actionData.errors?.general) {
        toast.error(actionData.errors.general);
      }
    }
  }, [actionData]);

  const onSubmit = (data: PostFormSchema) => {
    const formData = new FormData();
    formData.append("title", data.title);
    formData.append("content", data.content);
    data.category_ids.forEach((id) => {
      formData.append("categoryIds", id);
    });

    submit(formData, {
      method: "post",
    });
  };

  return (
    <div className="container w-full mx-auto">
      <div className="flex-1">
        <Form {...form}>
          <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 p-5">
            <div className="flex flex-row justify-between items-center">
              <h3 className="text-2xl font-bold">Add New Post</h3>

              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Creating..." : "Save"}
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
                    <Input disabled={field.disabled} {...field} />
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
                      disabled={field.disabled}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </ReactRouterForm>
        </Form>
      </div>
    </div>
  );
}
