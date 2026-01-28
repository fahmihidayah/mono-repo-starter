import z, { email } from "zod";

export const loginFormSchema = z.object({
    email : z.string().min(1, "Email is required").email("Invalid email address"),
    password : z.string().min(1, "Password is required").min(8, "Password must be at least 8 characters long"),
})

export type LoginFormData = z.infer<typeof loginFormSchema>;