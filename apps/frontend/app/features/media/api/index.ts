import type { Media, MediaFormSchema } from "~/features/media/types";
import { ApiClient } from "../../../lib/api";

export const mediaApi = new ApiClient<Media, MediaFormSchema>({
  resource: "api/media",
  defaultToken: "",
});
