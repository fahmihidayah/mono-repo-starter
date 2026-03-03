import { useResourceList } from "~/lib/hooks";

export function usePostList() {
  return useResourceList({
    successMessage: "Post created successfully!",
  });
}
