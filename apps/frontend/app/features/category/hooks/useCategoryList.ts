import { useSearchParams, useActionData, useSubmit } from "react-router";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export function useCategoryList() {
  const [searchParams, setSearchParams] = useSearchParams();
  const actionData = useActionData<{ success: boolean; message: string }>();
  const submit = useSubmit();
  const [searchValue, setSearchValue] = useState(searchParams.get("search") || "");

  useEffect(() => {
    if (actionData) {
      if (actionData.success) {
        toast.success(actionData.message || "Action completed successfully!");
      } else {
        toast.error(actionData.message || "Action failed");
      }
    }
  }, [actionData]);

  const handleSearch = (value: string) => {
    setSearchValue(value);
    const params = new URLSearchParams(searchParams);
    if (value) {
      params.set("search", value);
    } else {
      params.delete("search");
    }
    params.set("page", "1");
    setSearchParams(params);
  };

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams);
    params.set("page", newPage.toString());
    setSearchParams(params);
  };

  const handleBulkDelete = (selectedIds: string[]) => {
    const formData = new FormData();
    formData.append("intent", "bulkDelete");
    selectedIds.forEach((id) => {
      formData.append("ids[]", id);
    });
    submit(formData, { method: "post" });
  };

  return {
    searchValue,
    handleSearch,
    handlePageChange,
    handleBulkDelete,
  };
}
