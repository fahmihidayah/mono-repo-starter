import type { User } from "../types";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";

interface UserListProps {
  users: User[];
  pagination: {
    page: number;
    totalPages: number;
    totalDocs: number;
  };
  searchValue: string;
  onSearchChange: (value: string) => void;
  onPageChange: (page: number) => void;
  onBulkDelete: (ids: string[]) => void;
}

export function UserList({
  users,
  pagination,
  searchValue,
  onSearchChange,
  onPageChange,
  onBulkDelete,
}: UserListProps) {
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

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Users"
          description="Manage and organize your users efficiently"
          addButtonText="Add User"
          addButtonLink="/dashboard/users/add"
        />

        <DataTable
          data={users}
          columns={columns}
          enableRowSelection={true}
          onBulkDelete={onBulkDelete}
          getRowId={(row) => row.id}
          searchPlaceholder="Search users..."
          searchValue={searchValue}
          onSearchChange={onSearchChange}
          emptyMessage="No users found."
          totalPages={pagination.totalPages}
          manualPagination
        />

        {pagination.totalPages > 1 && (
          <TablePagination
            href="/dashboard/users"
            currentPage={pagination.page}
            totalPages={pagination.totalPages}
            onPageChange={onPageChange}
          />
        )}
      </div>
    </div>
  );
}
