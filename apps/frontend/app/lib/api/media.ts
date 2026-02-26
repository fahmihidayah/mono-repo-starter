import type { Media } from "~/features/media/type";
import { ApiClient } from "../api";

export const mediaApi = new ApiClient<Media>({
  resource: "api/media",
  defaultToken: "",
});
