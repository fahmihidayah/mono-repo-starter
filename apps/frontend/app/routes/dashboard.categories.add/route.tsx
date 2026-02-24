import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {
  Form as ReactRouterForm,
  redirect,
  useSubmit,
  useActionData,
  useNavigation,
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
import { categoryApi } from "~/lib/api/category";
import type { ActionData } from "~/types";

const CategoryFormSchema = z.object({
  title: z.string().min(1, "Title is required").max(100, "Title must be less than 100 characters"),
});

type CategoryFormData = z.infer<typeof CategoryFormSchema>;

export async function action({ request }: { request: Request }) {
  const formData = await request.formData();
  const data = Object.fromEntries(formData);
  const validation = CategoryFormSchema.safeParse(data);

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
    await categoryApi.create({
      request,
      data: data,
    });
    return redirect("/dashboard/categories?success=created");
  } catch (error) {
    return {
      success: false,
      errors: {
        general: "Failed to create category. Please try again.",
      },
    } as ActionData;
  }
}

export default function DashboardAddCategoryPage() {
  const form = useForm<CategoryFormData>({
    resolver: zodResolver(CategoryFormSchema),
    defaultValues: {
      title: "",
    },
  });

  const submit = useSubmit();
  const actionData = useActionData<ActionData>();
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

  // Show toast on error
  useEffect(() => {
    if (actionData && !actionData.success) {
      if (actionData.errors?.general) {
        toast.error(actionData.errors.general);
      }
    }
  }, [actionData]);

  const onSubmit = (data: CategoryFormData) => {
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
          <CardTitle className="text-2xl">Add New Category</CardTitle>
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

              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Creating..." : "Submit"}
              </Button>
            </ReactRouterForm>
          </Form>
        </CardContent>
      </Card>
      {/* Form for adding a new category */}
      {/* You can use the CategoryFormSchema for form validation */}
      {/* Example form: */}
      {/* <form method="post">
                <input type="text" name="title" placeholder="Category Title" required />
                <button type="submit">Create Category</button>
            </form> */}
    </div>
  );
}
