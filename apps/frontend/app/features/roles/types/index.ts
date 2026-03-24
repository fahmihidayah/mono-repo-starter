import z from "zod";

export type Role = {
  id: string;
  name: string | null;
  created_at: Date | null;
  updated_at: Date | null;
};

export const roleFormSchema = z.object({
  name: z.string().min(3, "Name is required and must be at least 3 characters").max(255, "Name must be less than 255 characters"),
});

export type RoleFormSchema = z.infer<typeof roleFormSchema>;
