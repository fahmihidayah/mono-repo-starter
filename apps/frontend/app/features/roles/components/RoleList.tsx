import type { Role } from "../types";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";

interface RoleListProps {
  roles: Role[];
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

export function RoleList({
  roles,
  pagination,
  searchValue,
  onSearchChange,
  onPageChange,
  onBulkDelete,
}: RoleListProps) {
  const columns = createColumn<Role>({
    columnConfig: [
      {
        type: "text",
        accessorKey: "name",
        header: "Name",
        fallback: "No name",
        isBold: true,
      },
      {
        type: "date",
        accessorKey: "created_at",
        header: "Created",
      },
      {
        type: "date",
        accessorKey: "updated_at",
        header: "Updated",
      },
    ],
    actionColumnConfig: {
      getItemId: (role) => role.id,
      editLink: (role) => `/dashboard/roles/${role.id}`,
      deleteFormAction: "/dashboard/roles",
    },
  });

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Roles"
          description="Manage and organize your roles efficiently"
          addButtonText="Add Role"
          addButtonLink="/dashboard/roles/add"
        />

        <DataTable
          data={roles}
          columns={columns}
          enableRowSelection={true}
          onBulkDelete={onBulkDelete}
          getRowId={(row) => row.id}
          searchPlaceholder="Search roles..."
          searchValue={searchValue}
          onSearchChange={onSearchChange}
          emptyMessage="No roles found."
          totalPages={pagination.totalPages}
          manualPagination
        />

        {pagination.totalPages > 1 && (
          <TablePagination
            href="/dashboard/roles"
            currentPage={pagination.page}
            totalPages={pagination.totalPages}
            onPageChange={onPageChange}
          />
        )}
      </div>
    </div>
  );
}
