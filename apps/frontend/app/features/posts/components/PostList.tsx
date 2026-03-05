import type { Post } from "../types";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";

interface PostListProps {
  posts: Post[];
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

export function PostList({
  posts,
  pagination,
  searchValue,
  onSearchChange,
  onPageChange,
  onBulkDelete,
}: PostListProps) {
  const columns = createColumn<Post>({
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
      getItemId: (post) => post.id,
      editLink: (post) => `/dashboard/posts/${post.id}`,
      deleteFormAction: "/dashboard/posts",
    },
  });

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Posts"
          description="Manage and organize your posts efficiently"
          addButtonText="Add Post"
          addButtonLink="/dashboard/posts/add"
        />

        <DataTable
          data={posts}
          columns={columns}
          enableRowSelection={true}
          onBulkDelete={onBulkDelete}
          getRowId={(row) => row.id || ""}
          searchPlaceholder="Search posts..."
          searchValue={searchValue}
          onSearchChange={onSearchChange}
          emptyMessage="No posts found."
          totalPages={pagination.totalPages}
          manualPagination
        />

        {pagination.totalPages > 1 && (
          <TablePagination
            href="/dashboard/posts"
            currentPage={pagination.page}
            totalPages={pagination.totalPages}
            onPageChange={onPageChange}
          />
        )}
      </div>
    </div>
  );
}
