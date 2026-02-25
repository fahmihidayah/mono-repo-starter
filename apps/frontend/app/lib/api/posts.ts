import type { Post } from "~/features/posts/type";
import { ApiClient } from "../api";

export const postApi = new ApiClient<Partial<Post>>({
  resource: "api/posts",
  defaultToken: "",
});
