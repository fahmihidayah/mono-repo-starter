import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import { z } from "zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "~/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form as ReactRouterForm, redirect, useLoaderData, useSubmit } from "react-router";
import { taskRepository } from "~/features/tasks/task-repository";
import type { Route } from "./+types/dashboard.tasks.$id";


// Zod schema
const taskSchema = z.object({
    title: z.string().min(1, "Title is required").max(100, "Title must be less than 100 characters"),
    description: z.string().min(10, "Description must be at least 10 characters").max(500, "Description must be less than 500 characters"),
});

type TaskFormData = z.infer<typeof taskSchema>;


export async function loader({ params }: Route.LoaderArgs) {
    // Use repository's findById method
    const task = await taskRepository.findById(params.id);

    if (!task) {
        throw new Response("Task not found", { status: 404 });
    }

    return task;
}

// Server action - Update task using TaskRepository
export async function action({ request, params }: Route.ActionArgs) {
    const formData = await request.formData();

    try {
        const title = formData.get('title')?.toString();
        const description = formData.get("description")?.toString();

        // Use repository's update method
        await taskRepository.update(params.id, {
            title,
            description,
            updatedAt: new Date(),
        });

        return redirect("/dashboard/tasks");
    } catch (error) {
        return {
            error: "Failed to update task. Please try again.",
        };
    }
}

export default function EditTaskPage() {
    const loaderData = useLoaderData<typeof loader>();
    const form = useForm<TaskFormData>({
        resolver: zodResolver(taskSchema),
        defaultValues: {
            title: loaderData.title || "",
            description: loaderData.description || ""
        }
    })

    const submit = useSubmit();

    const onSubmit = (data: TaskFormData) => {
        const formData = new FormData();
        formData.append("title", data.title);
        formData.append("description", data.description);

        // Sends to your `action()` on this route
        submit(formData, {
            method: "post",
        });
    };
    return (
        <div className="container w-full mx-auto p-5">
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl">Edit Task</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                    <Form {...form}>
                        <ReactRouterForm
                            method="post"
                            onSubmit={form.handleSubmit(onSubmit)}
                            className="space-y-2"
                        >
                            <FormField
                                control={form.control}
                                name="title"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Title</FormLabel>
                                        <FormControl>
                                            <Input
                                                type="text"
                                                disabled={field.disabled}
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="description"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Description</FormLabel>
                                        <FormControl>
                                            <Textarea
                                                disabled={field.disabled}
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <Button type="submit">
                                Update Task
                            </Button>
                        </ReactRouterForm>
                    </Form>
                </CardContent>
            </Card>
        </div>
    )
}