import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import z from "zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "~/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form as ReactRouterForm, redirect, useLoaderData, useSubmit } from "react-router";
import { userApi } from "~/lib/api/users";
import type { ActionData } from "~/types";
import type { Route } from "../dashboard.users.$id/+types/route";

// Zod schema
const userSchema = z.object({
    name: z.string().min(1, "Name is required").max(100, "Name must be less than 100 characters"),
    email: z.string().email("Invalid email address"),
});

type UserFormData = z.infer<typeof userSchema>;

// Server action
export async function action({ request, params }: Route.ActionArgs) {
    const id = params.id;
    const formData = await request.formData();
    const data = Object.fromEntries(formData);
    const validation = userSchema.safeParse(data);
    console.log(validation);

    if (!validation.success) {
        const fieldErrors = validation.error.flatten((issue) => issue.message).fieldErrors;
        return {
            success: false,
            errors: {
                ...fieldErrors,
            }
        } as ActionData;
    }
    try {
        const response = await userApi.update({
            request,
            data: data,
            id
        });
        console.log(response);
        return redirect("/dashboard/users");
    } catch (error) {
        return {
            error: "Failed to create user. Please try again.",
        };
    }
}

export async function loader({ request, params }: Route.LoaderArgs) {
    const id = params.id;
    const user = await userApi.getById({
        request,
        id
    });
    return { user };
}

export default function AddUserPage() {
    const { user } = useLoaderData<typeof loader>();
    const form = useForm<UserFormData>({
        resolver: zodResolver(userSchema),
        defaultValues: {
            name: user.data?.name,
            email: user.data?.email
        }
    })

    const submit = useSubmit();

    const onSubmit = (data: UserFormData) => {
        const formData = new FormData();
        formData.append("name", data.name);
        formData.append("email", data.email);

        submit(formData, {
            method: "post",
        });
    };

    return <div className="container w-full mx-auto p-5">
        <Card>
            <CardHeader>
                <CardTitle className="text-2xl">Add New User</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
                <Form {...form}>
                    <ReactRouterForm onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
                        <FormField control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Name</FormLabel>
                                    <FormControl>
                                        <Input
                                            type="text"
                                            disabled={field.disabled}
                                            {...field}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )} />
                        <FormField control={form.control}
                            name="email"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Email</FormLabel>
                                    <FormControl>
                                        <Input
                                            type="email"
                                            disabled={field.disabled}
                                            {...field}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )} />

                        <Button>
                            Submit
                        </Button>
                    </ReactRouterForm>
                </Form>
            </CardContent>
        </Card>
    </div>
}
