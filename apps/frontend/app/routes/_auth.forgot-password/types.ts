import { z } from "zod";

export const forgotPasswordFormSchema = z.object({
  email: z.string().min(1, "Email is required").email("Invalid email address"),
});

export type ForgotPasswordFormData = z.infer<typeof forgotPasswordFormSchema>;
