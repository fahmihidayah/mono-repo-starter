import z from "zod";

export type Media = {
  id: string;
  alt: string | null;
  url: string | null;
  path: string | null;
  file_name: string | null;
  mime_type: string | null;
  file_size: number | null;
  width: number | null;
  height: number | null;
  created_at: number | null;
  updated_at: number | null;
};

export const mediaFormSchema = z.object({
  alt: z.string().max(255, "Alt text must be less than 255 characters").nullable(),
  url: z.string().max(255, "URL must be less than 255 characters").nullable(),
});

export type MediaFormSchema = z.infer<typeof mediaFormSchema>;
