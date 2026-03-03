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
import { categoryApi } from "~/features/category/api";
import type { ActionData } from "~/types";
import type { Route } from "./+types/route";
import { categoryFormSchema, type CategoryFormSchema } from "~/features/category/types";

export async function loader({ request, params }: Route.LoaderArgs) {
  const id = params.id;
  const category = await categoryApi.getById({
    request,
    id,
  });
  return { category };
}

export async function action({ request, params }: Route.ActionArgs) {
  const id = params.id;
  const formData = await request.formData();
  const data = Object.fromEntries(formData);
  const validation = categoryFormSchema.safeParse(data);

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
    await categoryApi.update({
      request,
      data: data,
      id,
    });
    return redirect("/dashboard/categories");
  } catch (error) {
    return {
      success: false,
      errors: {
        general: "Failed to update category. Please try again.",
      },
    } as ActionData;
  }
}

export default function EditCategoryPage() {
  const { category } = useLoaderData<typeof loader>();
  const form = useForm<CategoryFormSchema>({
    resolver: zodResolver(categoryFormSchema),
    defaultValues: {
      title: category.data?.title || "",
    },
  });

  const submit = useSubmit();

  const onSubmit = (data: CategoryFormSchema) => {
    const formData = new FormData();
    formData.append("title", data.title);

    submit(formData, {
      method: "post",
    });
  };

  return (
    <div className="container w-full mx-auto p-5">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Edit Category</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
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

              <Button type="submit">Update</Button>
            </ReactRouterForm>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
