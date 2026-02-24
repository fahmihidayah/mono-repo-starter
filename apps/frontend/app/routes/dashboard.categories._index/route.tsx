import type { Category } from "~/features/category/type";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";
import { useLoaderData, useSearchParams, useActionData } from "react-router";
import { categoryApi } from "~/lib/api/category";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export async function loader({ request }: { request: Request }) {
  const url = new URL(request.url);
  const page = Number.parseInt(url.searchParams.get("page") || "1");
  const pageSize = Number.parseInt(url.searchParams.get("limit") || "10");
  const search = url.searchParams.get("search") || "";

  // Use repository's findManyPaginated method
  const result = await categoryApi.getAll({
    request,
    page,
    pageSize,
    search: { title: search },
  });

  return result;
}

// Action - Handle delete and update operations
export async function action({ request }: { request: Request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const id = formData.get("id")?.toString();

  try {
    if (intent === "delete" && id) {
      await categoryApi.deleteById({
        request,
        id,
      });
      return { success: true, message: "Category deleted successfully" };
    }
    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}

export default function CategoriesDashboardPage() {
  const loaderData = useLoaderData<typeof loader>();
  const actionData = useActionData<typeof action>();

  const [searchParams, setSearchParams] = useSearchParams();

  // State
  const [searchValue, setSearchValue] = useState(searchParams.get("search") || "");

  // Show toast on delete action
  useEffect(() => {
    if (actionData) {
      if (actionData.success) {
        toast.success(actionData.message || "Category deleted successfully!");
      } else {
        toast.error(actionData.message || "Failed to delete category");
      }
    }
  }, [actionData]);

  const columns = createColumn<Category>({
    columnConfig: [
      {
        type: "text",
        accessorKey: "slug",
        header: "Slug",
        fallback: "No slug",
        isBold: true,
      },
      {
        type: "text",
        accessorKey: "title",
        header: "Title",
        fallback: "No Title",
      },
      {
        type: "date",
        accessorKey: "createdAt",
        header: "Created",
      },
      {
        type: "date",
        accessorKey: "updatedAt",
        header: "Updated",
      },
    ],
    actionColumnConfig: {
      getItemId: (category) => category.id,
      editLink: (category) => `/dashboard/categories/${category.id}`,
      deleteFormAction: "/dashboard/categories",
    },
  });

  const handleSearch = (value: string) => {
    setSearchValue(value);
    const params = new URLSearchParams(searchParams);
    if (value) {
      params.set("search", value);
    } else {
      params.delete("search");
    }
    params.set("page", "1"); // Reset to first page
    setSearchParams(params);
  };

  // Handle page change
  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams);
    params.set("page", newPage.toString());
    setSearchParams(params);
  };

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        {/* Page Header */}
        <PageHeader
          title="Categories"
          description="Manage and organize your categories efficiently"
          addButtonText="Add Category"
          addButtonLink="/dashboard/categories/add"
        />

        {/* Data Table */}
        <DataTable
          title="All Categorie"
          description={`${loaderData.pagination?.totalDocs} categories${loaderData.pagination?.totalDocs !== 1 ? "" : ""} total`}
          data={loaderData.data ?? []}
          columns={columns}
          searchPlaceholder="Search users..."
          searchValue={searchValue}
          onSearchChange={handleSearch}
          emptyMessage="No tasks found."
          totalPages={loaderData.pagination?.totalPages ?? 0}
          manualPagination
        />

        {/* Table Pagination */}
        {(loaderData.pagination?.totalPages ?? 0) > 1 && (
          <TablePagination
            href={"/dashboard/categories"}
            currentPage={loaderData.pagination?.page ?? 1}
            totalPages={loaderData.pagination?.totalPages ?? 0}
            onPageChange={handlePageChange}
          />
        )}
      </div>
    </div>
  );
}
