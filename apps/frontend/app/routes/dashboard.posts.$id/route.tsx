import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Form as ReactRouterForm, redirect, useLoaderData, useSubmit } from "react-router";
import { z } from "zod";
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
import { postApi } from "~/lib/api/posts";
import { categoryApi } from "~/lib/api/category";
import type { ActionData } from "~/types";
import type { Route } from "./+types/route";

const PostFormSchema = z.object({
  title: z.string().min(1, "Title is required").max(255, "Title must be less than 255 characters"),
  content: z.string().min(1, "Content is required"),
  categoryIds: z.array(z.string()).min(1, "Select at least one category"),
});

type PostFormData = z.infer<typeof PostFormSchema>;

export async function loader({ request, params }: Route.LoaderArgs) {
  const id = params.id;
  const [post, categoriesResponse] = await Promise.all([
    postApi.getById({ request, id }),
    categoryApi.getAll({ request, page: 1, pageSize: 100 }),
  ]);

  return {
    post: post.data,
    categories: categoriesResponse.data || [],
  };
}

export async function action({ request, params }: Route.ActionArgs) {
  const id = params.id;
  const formData = await request.formData();
  const categoryIds = formData.getAll("categoryIds");
  const data = {
    title: formData.get("title"),
    content: formData.get("content"),
    categoryIds: categoryIds,
  };

  const validation = PostFormSchema.safeParse(data);

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
    await postApi.update({
      request,
      data: {
        title: validation.data.title,
        content: validation.data.content,
        category_ids: validation.data.categoryIds,
      },
      id,
    });
    return redirect("/dashboard/posts");
  } catch (error) {
    return {
      success: false,
      errors: {
        general: "Failed to update post. Please try again.",
      },
    } as ActionData;
  }
}

export default function EditPostPage() {
  const { post, categories } = useLoaderData<typeof loader>();
  const form = useForm<PostFormData>({
    resolver: zodResolver(PostFormSchema),
    defaultValues: {
      title: post?.title || "",
      content: post?.content || "",
      categoryIds: post?.categories?.map((cat) => cat.id) || [],
    },
  });

  const submit = useSubmit();

  const onSubmit = (data: PostFormData) => {
    const formData = new FormData();
    formData.append("title", data.title);
    formData.append("content", data.content);
    data.categoryIds.forEach((id) => {
      formData.append("categoryIds", id);
    });

    submit(formData, {
      method: "post",
    });
  };

  return (
    <div className="container w-full mx-auto p-5">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Edit Post</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                      <Textarea disabled={field.disabled} rows={10} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="categoryIds"
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

              <Button type="submit">Update</Button>
            </ReactRouterForm>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
