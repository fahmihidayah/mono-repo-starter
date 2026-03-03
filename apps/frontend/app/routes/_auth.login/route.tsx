import { Link, Form, redirect, useActionData, useNavigation, useSubmit } from "react-router";
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
import { PasswordInput } from "~/components/ui/password-input";
import { getSession, commitSession } from "~/session.server";
// import { authenticateUser, calculateSessionMaxAge } from "~/lib/auth.server";
import type { Route } from "../_auth.login/+types/route";
import { useForm } from "react-hook-form";
import { loginFormSchema, type LoginFormData } from "./types";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form as FormComponent,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Lock, Mail } from "lucide-react";
import type { ActionData } from "~/types";
import { authApi } from "~/features/users/api/auth";
import { calculateSessionMaxAge } from "~/lib/utils.server";

// Meta function for SEO
export function meta() {
  return [
    { title: "Login - Starter App" },
    { name: "description", content: "Sign in to your account" },
  ];
}

// export type ActionData = {
//   success: boolean;
//   errors: { [key: string]: string | undefined }
// };

// Server action for login
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const data = Object.fromEntries(formData);

  // Validate form data with Zod
  const validation = loginFormSchema.safeParse(data);

  if (!validation.success) {
    const fieldErrors = validation.error.flatten((issue) => issue.message).fieldErrors;
    return {
      success: false,
      errors: {
        ...fieldErrors,
      },
    } as ActionData;
  }

  // Authenticate with Go API
  const result = await authApi.login(validation.data);

  if (result.code !== 200 || !result.data) {
    return {
      success: false,
      errors: {
        general: result.message || "Invalid email or password. Please try again.",
      },
    } as ActionData;
  }

  // Create session with auth data
  const session = await getSession(request.headers.get("Cookie"));
  const { token, id, name, email, exp } = result.data;

  session.set("userId", id);
  session.set("token", token);
  session.set("userName", name);
  session.set("userEmail", email);

  // Calculate session expiration from token
  const maxAge = calculateSessionMaxAge(exp);

  // Commit session and redirect
  return redirect("/dashboard", {
    headers: {
      "Set-Cookie": await commitSession(session, { maxAge }),
    },
  });
}

// Login component
export default function Login() {
  const actionData = useActionData<typeof action>();
  const navigation = useNavigation();
  const submit = useSubmit();
  const isSubmitting = navigation.state === "submitting";

  // Initialize React Hook Form with Zod validation
  const form = useForm<LoginFormData>({
    resolver: zodResolver(loginFormSchema),
    mode: "onBlur",
  });

  // Submit handler - delegates to React Router action
  const onSubmit = (data: LoginFormData) => {
    submit(data, { method: "post" });
  };
  return (
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-background via-background to-muted/20 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">Welcome back</CardTitle>
          <CardDescription className="text-center">
            Enter your credentials to access your account
          </CardDescription>
        </CardHeader>
        <FormComponent {...form}>
          <Form method="post" onSubmit={form.handleSubmit(onSubmit)}>
            <CardContent className="flex flex-col gap-4">
              {/* General error message */}
              {actionData?.errors.general && (
                <div className="bg-destructive/10 text-destructive text-sm p-3 rounded-md border border-destructive/20">
                  {actionData.errors.general}
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

              {/* Password field with show/hide toggle */}
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground z-10 pointer-events-none" />
                        <PasswordInput
                          className="pl-10"
                          autoComplete="current-password"
                          placeholder="••••••••"
                          required
                          disabled={isSubmitting}
                          aria-invalid={actionData?.errors["password"] ? "true" : undefined}
                          aria-describedby={
                            actionData?.errors?.password ? "password-error" : undefined
                          }
                          {...field}
                          value={(field.value as string) ?? ""}
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
            <CardFooter className="flex flex-col pt-5 space-y-4">
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? "Signing in..." : "Sign in"}
              </Button>

              <div className="text-sm text-center text-muted-foreground">
                <Link to="/forgot-password" className="text-primary hover:underline font-medium">
                  Forgot password?
                </Link>
              </div>

              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-card px-2 text-muted-foreground">Or</span>
                </div>
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
