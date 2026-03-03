import type { Category, CategoryFormSchema } from "~/features/category/types";
import { ApiClient } from "../../../lib/api";

export const categoryApi = new ApiClient<Category, CategoryFormSchema>({
  resource: "api/categories",
  defaultToken: "",
});
