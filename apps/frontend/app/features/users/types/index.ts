import z from "zod";

export interface User {
  id: string;
  name: string;
  email: string;
  password?: string;
  emailVerified: boolean;
  image?: string;
  createdAt: Date;
  updatedAt: Date;
}

// Zod schema
export const userFormSchema = z.object({
  name: z.string().min(1, "Name is required").max(100, "Name must be less than 100 characters"),
  email: z.email("Invalid email address"),
  password: z.string().min(8, "Password must be at least 8 characters long"),
});

export type UserFormSchema = z.infer<typeof userFormSchema>;
