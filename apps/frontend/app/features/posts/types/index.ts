import z from "zod";
import type { Category } from "../../category/types";

export type Post = {
  id: string;
  slug: string | null;
  title: string | null;
  content: string | null;
  categories: Category[] | null;
  user_id: string | null;
  created_at: Date | null;
  updated_at: Date | null;
};

export const postFormSchema = z.object({
  title: z.string().min(1, "Title is required").max(255, "Title must be less than 255 characters"),
  content: z.string().min(1, "Content is required"),
  category_ids: z.array(z.string()).min(1, "Select at least one category"),
});

export type PostFormSchema = z.infer<typeof postFormSchema>;
