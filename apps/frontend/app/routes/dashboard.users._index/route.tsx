import { useEffect, useState } from "react";
import { useLoaderData, useSearchParams, useActionData, useSubmit } from "react-router";
import { PageHeader, DataTable, TablePagination } from "~/components/layout/table/table-list";
import createColumn from "~/components/layout/table/column/create-column";
import type { User } from "~/features/users/types";
import { userApi } from "~/features/users/api/users";
import { toast } from "sonner";

// Loader - Fetch users with pagination and search using UserRepository
export async function loader({ request }: { request: Request }) {
  const url = new URL(request.url);
  const page = Number.parseInt(url.searchParams.get("page") || "1");
  const pageSize = Number.parseInt(url.searchParams.get("limit") || "10");
  const search = url.searchParams.get("search") || "";

  // Use repository's findManyPaginated method
  const result = await userApi.getAll({
    request,
    page,
    pageSize,
    search: { email: search },
  });

  return result;
}

// Action - Handle delete and update operations
export async function action({ request }: { request: Request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const id = formData.get("id")?.toString();
  const ids = formData.getAll("ids[]");

  try {
    if (intent === "delete" && id) {
      await userApi.deleteById({
        request,
        id,
      });
      return { success: true, message: "User deleted successfully" };
    }

    if (intent === "bulkDelete" && ids.length > 0) {
      await userApi.deleteMany({
        request,
        ids: ids.map((id) => id.toString()),
      });
      return {
        success: true,
        message: `${ids.length} user${ids.length !== 1 ? "s" : ""} deleted successfully`,
      };
    }

    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}

export function meta() {
  return [{ title: "User - Dashboard" }, { name: "description", content: "Manage your tasks" }];
}

export default function DashboardUsersPage() {
  const loaderData = useLoaderData<typeof loader>();
  const actionData = useActionData<typeof action>();
  const submit = useSubmit();
  const [searchParams, setSearchParams] = useSearchParams();

  // State
  const [searchValue, setSearchValue] = useState(searchParams.get("search") || "");

  // Show toast on delete action
  useEffect(() => {
    if (actionData) {
      if (actionData.success) {
        toast.success(actionData.message || "User deleted successfully!");
      } else {
        toast.error(actionData.message || "Failed to delete user");
      }
    }
  }, [actionData]);

  // Table columns
  const columns = createColumn<User>({
    columnConfig: [
      {
        type: "text",
        accessorKey: "email",
        header: "Email",
        fallback: "No email",
        isBold: true,
      },
      {
        type: "text",
        accessorKey: "name",
        header: "Name",
        fallback: "No Name",
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
      getItemId: (user) => user.id,
      editLink: (user) => `/dashboard/users/${user.id}`,
      deleteFormAction: "/dashboard/users",
    },
  });

  // Handle search
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

  const handleBulkDelete = (selectedIds: string[]) => {
    const formData = new FormData();
    formData.append("intent", "bulkDelete");
    selectedIds.forEach((id) => {
      formData.append("ids[]", id);
    });
    submit(formData, { method: "post" });
  };

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        {/* Page Header */}
        <PageHeader
          title="Users"
          description="Manage and organize your users efficiently"
          addButtonText="Add User"
          addButtonLink="/dashboard/users/add"
        />

        {/* Data Table */}
        <DataTable
          title="All Users"
          description={`${loaderData.pagination?.totalDocs} user${loaderData.pagination?.totalDocs !== 1 ? "s" : ""} total`}
          data={loaderData.data ?? []}
          columns={columns}
          enableRowSelection={true}
          onBulkDelete={handleBulkDelete}
          getRowId={(row) => row.id}
          searchPlaceholder="Search users..."
          searchValue={searchValue}
          onSearchChange={handleSearch}
          emptyMessage="No users found."
          totalPages={loaderData.pagination?.totalPages ?? 0}
          manualPagination
        />

        {/* Table Pagination */}
        {(loaderData.pagination?.totalPages ?? 0) > 1 && (
          <TablePagination
            href={"/dashboard/users"}
            currentPage={loaderData.pagination?.page ?? 1}
            totalPages={loaderData.pagination?.totalPages ?? 0}
            onPageChange={handlePageChange}
          />
        )}
      </div>
    </div>
  );
}
