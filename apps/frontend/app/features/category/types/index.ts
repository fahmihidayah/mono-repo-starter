import z from "zod";

export type Category = {
  id: string;
  slug: string | null;
  title: string | null;
  createdAt: Date | null;
  updatedAt: Date | null;
};

export const categoryFormSchema = z.object({
  title: z.string().min(1, "Title is required").max(100, "Title must be less than 100 characters"),
});

export type CategoryFormSchema = z.infer<typeof categoryFormSchema>;
