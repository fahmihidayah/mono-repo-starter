import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import z from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form as ReactRouterForm, redirect, useSubmit } from "react-router";
import { userApi } from "~/features/users/api/users";
import type { ActionData } from "~/types";
import { userFormSchema, type UserFormSchema } from "~/features/users/types";

// Server action
export async function action({ request }: { request: Request }) {
  const formData = await request.formData();
  const data = Object.fromEntries(formData);
  const validation = userFormSchema.safeParse(data);
  console.log(validation);

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
    const response = await userApi.create({
      request,
      data: data,
    });
    console.log(response);
    return redirect("/dashboard/users");
  } catch (error) {
    return {
      error: "Failed to create user. Please try again.",
    };
  }
}

export default function AddUserPage() {
  const form = useForm<UserFormSchema>({
    resolver: zodResolver(userFormSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
    },
  });

  const submit = useSubmit();

  const onSubmit = (data: UserFormSchema) => {
    const formData = new FormData();
    formData.append("name", data.name);
    formData.append("email", data.email);
    formData.append("password", data.password);

    submit(formData, {
      method: "post",
    });
  };

  return (
    <div className="container w-full mx-auto p-5">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Add New User</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <Form {...form}>
            <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input type="text" disabled={field.disabled} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input type="email" disabled={field.disabled} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input type="password" disabled={field.disabled} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button>Submit</Button>
            </ReactRouterForm>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
