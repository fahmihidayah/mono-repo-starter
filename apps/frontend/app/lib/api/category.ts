import type { Category } from "~/features/category/type";
import { ApiClient } from "../api";

export const categoryApi = new ApiClient<Category>({
    resource: "api/categories",
    defaultToken: ""
});