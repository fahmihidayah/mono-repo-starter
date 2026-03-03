import { useResourceList } from "~/lib/hooks";

export function useCategoryList() {
  return useResourceList({
    successMessage: "Category created successfully!",
  });
}
