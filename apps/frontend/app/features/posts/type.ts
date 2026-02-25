import type { Category } from "../category/type";

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
