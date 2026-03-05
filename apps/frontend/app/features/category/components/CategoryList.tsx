import type { Category } from "../types";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";

interface CategoryListProps {
  categories: Category[];
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

export function CategoryList({
  categories,
  pagination,
  searchValue,
  onSearchChange,
  onPageChange,
  onBulkDelete,
}: CategoryListProps) {
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

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Categories"
          description="Manage and organize your categories efficiently"
          addButtonText="Add Category"
          addButtonLink="/dashboard/categories/add"
        />

        <DataTable
          data={categories}
          columns={columns}
          enableRowSelection={true}
          onBulkDelete={onBulkDelete}
          getRowId={(row) => row.id}
          searchPlaceholder="Search categories..."
          searchValue={searchValue}
          onSearchChange={onSearchChange}
          emptyMessage="No categories found."
          totalPages={pagination.totalPages}
          manualPagination
        />

        {pagination.totalPages > 1 && (
          <TablePagination
            href="/dashboard/categories"
            currentPage={pagination.page}
            totalPages={pagination.totalPages}
            onPageChange={onPageChange}
          />
        )}
      </div>
    </div>
  );
}
