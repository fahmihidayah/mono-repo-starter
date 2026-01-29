import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useSubmit, useActionData, useNavigation } from "react-router";
import { useEffect } from "react";
import { toast } from "sonner";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { User, Mail } from "lucide-react";
import type { ActionData } from "~/types";

const profileFormSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters").max(100, "Name must not exceed 100 characters"),
  email: z.string().email("Please enter a valid email address"),
  actionType: z.literal("updateProfile"),
});

export type ProfileFormData = z.infer<typeof profileFormSchema>;

interface UpdateProfileFormProps {
  defaultValues: {
    name: string;
    email: string;
  };
}

export function UpdateProfileForm({ defaultValues }: UpdateProfileFormProps) {
  const submit = useSubmit();
  const actionData = useActionData<ActionData>();
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting" &&
    navigation.formData?.get("actionType") === "updateProfile";

  const form = useForm<ProfileFormData>({
    resolver: zodResolver(profileFormSchema),
    defaultValues: {
      name: defaultValues.name,
      email: defaultValues.email,
      actionType: "updateProfile",
    },
    mode: "onBlur",
  });

  // Show toast on success or error
  useEffect(() => {
    if (actionData?.actionType === "updateProfile") {
      if (actionData.success) {
        toast.success("Profile updated successfully!");
      } else if (actionData.errors?.general) {
        toast.error(actionData.errors.general);
      }
    }
  }, [actionData]);

  const onSubmit = (data: ProfileFormData) => {
    submit(data, { method: "post" });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Profile Information</CardTitle>
        <CardDescription>
          Update your personal information and email address
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <input type="hidden" name="actionType" value="updateProfile" />

            {actionData?.errors?.general && actionData?.actionType === "updateProfile" && (
              <div className="bg-destructive/10 text-destructive text-sm p-3 rounded-md border border-destructive/20">
                {actionData.errors.general}
              </div>
            )}

            {actionData?.success && actionData?.actionType === "updateProfile" && (
              <div className="bg-green-50 text-green-700 text-sm p-3 rounded-md border border-green-200">
                Profile updated successfully!
              </div>
            )}

            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                      <Input
                        className="pl-10"
                        placeholder="Your name"
                        disabled={isSubmitting}
                        {...field}
                      />
                    </div>
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
                    <div className="relative">
                      <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                      <Input
                        className="pl-10"
                        type="email"
                        placeholder="your@email.com"
                        disabled={isSubmitting}
                        {...field}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Saving..." : "Save Changes"}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
