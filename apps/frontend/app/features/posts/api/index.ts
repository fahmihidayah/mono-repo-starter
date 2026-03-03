import type { Post, PostFormSchema } from "~/features/posts/types";
import { ApiClient } from "../../../lib/api";

export const postApi = new ApiClient<Post, PostFormSchema>({
  resource: "api/posts",
  defaultToken: "",
});
