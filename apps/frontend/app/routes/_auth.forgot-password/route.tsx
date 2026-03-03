import { Link, Form, useActionData, useNavigation, useSubmit } from "react-router";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import type { Route } from "../_auth.forgot-password/+types/route";
import { useForm } from "react-hook-form";
import { forgotPasswordFormSchema, type ForgotPasswordFormData } from "./types";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form as FormComponent,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Mail } from "lucide-react";
import type { ActionData } from "~/types";
import { authApi } from "~/features/users/api/auth";

// Meta function for SEO
export function meta() {
  return [
    { title: "Forgot Password - Starter App" },
    { name: "description", content: "Reset your password" },
  ];
}

// Server action for forgot password
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const data = Object.fromEntries(formData);

  // Validate form data with Zod
  const validation = forgotPasswordFormSchema.safeParse(data);

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
    const result = await authApi.forgotPassword(validation.data.email);

    if (result.code !== 200) {
      return {
        success: true,
        errors: {
          general: result.message || "Failed to send password reset email",
        },
      } as ActionData;
    }

    return {
      success: false,
      errors: {
        general: result.message || "Failed to send password reset email",
      },
    } as ActionData;
  } catch (error: any) {
    return {
      success: false,
      errors: {
        general: error.message || "Network error. Please try again.",
      },
    } as ActionData;
  }
}

// Forgot Password component
export default function ForgotPassword() {
  const actionData = useActionData<typeof action>();
  const navigation = useNavigation();
  const submit = useSubmit();
  const isSubmitting = navigation.state === "submitting";

  // Initialize React Hook Form with Zod validation
  const form = useForm<ForgotPasswordFormData>({
    resolver: zodResolver(forgotPasswordFormSchema),
    mode: "onBlur",
  });

  // Submit handler - delegates to React Router action
  const onSubmit = (data: ForgotPasswordFormData) => {
    submit(data, { method: "post" });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-background via-background to-muted/20 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">Forgot Password</CardTitle>
          <CardDescription className="text-center">
            Enter your email address and we'll send you a link to reset your password
          </CardDescription>
        </CardHeader>
        <FormComponent {...form}>
          <Form method="post" onSubmit={form.handleSubmit(onSubmit)}>
            <CardContent className="flex flex-col gap-4">
              {/* General error message */}
              {actionData?.errors?.general && (
                <div className="bg-destructive/10 text-destructive text-sm p-3 rounded-md border border-destructive/20">
                  {actionData.errors.general}
                </div>
              )}

              {/* Success message */}
              {actionData?.success && (
                <div className="bg-green-50 text-green-700 text-sm p-3 rounded-md border border-green-200">
                  Password reset email sent successfully! Check your inbox for further instructions.
                </div>
              )}

              {/* Email field */}
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
                          autoComplete="email"
                          required
                          disabled={isSubmitting}
                          aria-invalid={actionData?.errors?.email ? "true" : undefined}
                          aria-describedby={actionData?.errors?.email ? "email-error" : undefined}
                          {...field}
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
            <CardFooter className="flex flex-col pt-5 space-y-5">
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? "Sending..." : "Send Reset Link"}
              </Button>
              <div className="text-sm text-center text-muted-foreground">
                Remember your password?{" "}
                <Link to="/login" className="text-primary hover:underline font-medium">
                  Sign in
                </Link>
              </div>
              <div className="text-sm text-center text-muted-foreground">
                Don't have an account?{" "}
                <Link to="/register" className="text-primary hover:underline font-medium">
                  Sign up
                </Link>
              </div>
            </CardFooter>
          </Form>
        </FormComponent>
      </Card>
    </div>
  );
}
