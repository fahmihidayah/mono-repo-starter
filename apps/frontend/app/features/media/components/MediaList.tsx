import type { Media } from "../types";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";

interface MediaListProps {
  media: Media[];
  pagination: {
    page: number;
    totalPages: number;
    totalDocs: number;
  };
  searchValue: string;
  onSearchChange: (value: string) => void;
  onPageChange: (page: number) => void;
  onBulkDelete: (selectedIds: string[]) => void;
}

export function MediaList({
  media,
  pagination,
  searchValue,
  onSearchChange,
  onPageChange,
  onBulkDelete,
}: MediaListProps) {
  const columns = createColumn<Media>({
    columnConfig: [
      {
        accessorKey: "url",
        type: "image",
        header: "Image",
        fallback: "/default-image.png",
      },
      {
        type: "text",
        accessorKey: "file_name",
        header: "File Name",
        fallback: "No file name",
        isBold: true,
      },
      {
        type: "url",
        accessorKey: "url",
        header: "URL",
        fallback: "No URL",
        maxLength: 50,
      },
      {
        type: "custom",
        accessorKey: "file_size",
        header: "Size",
        cell: ({ row }) => {
          const size = row.original.file_size;
          if (!size) return <span className="text-muted-foreground">-</span>;

          const formatSize = (bytes: number) => {
            if (bytes < 1024) return `${bytes} B`;
            if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
            return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
          };

          return <span className="text-sm">{formatSize(size)}</span>;
        },
      },
    ],
    actionColumnConfig: {
      getItemId: (media) => media.id,
      deleteFormAction: "/dashboard/media",
    },
  });

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Media"
          description="Manage and organize your media files"
          addButtonText="Upload Media"
          addButtonLink="/dashboard/media/add"
        />

        <DataTable
          data={media}
          columns={columns}
          searchPlaceholder="Search media files..."
          searchValue={searchValue}
          onSearchChange={onSearchChange}
          emptyMessage="No media files found."
          totalPages={pagination.totalPages}
          manualPagination
          enableRowSelection={true}
          onBulkDelete={onBulkDelete}
          getRowId={(row) => row.id}
        />

        {pagination.totalPages > 1 && (
          <TablePagination
            href="/dashboard/media"
            currentPage={pagination.page}
            totalPages={pagination.totalPages}
            onPageChange={onPageChange}
          />
        )}
      </div>
    </div>
  );
}
