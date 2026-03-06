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
  image: z
    .custom<File>((val) => val instanceof File, {
      message: "Image is required",
    })
    .refine((file) => file.size > 0, { message: "File cannot be empty" })
    .refine((file) => file.size <= 5 * 1024 * 1024, { message: "Max file size is 5MB" })
    .refine((file) => ["image/jpeg", "image/png", "image/webp"].includes(file.type), {
      message: "Only .jpg, .png, .webp formats are allowed",
    }),
  alt: z.string().optional(),
});

export type MediaFormSchema = z.infer<typeof mediaFormSchema>;
