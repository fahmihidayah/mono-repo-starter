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

export type MediaFormSchema = {
  alt?: string | null;
  url?: string | null;
};
