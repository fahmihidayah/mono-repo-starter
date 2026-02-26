import type { Media } from "~/features/media/type";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";
import { useLoaderData, useSearchParams, useActionData, useSubmit } from "react-router";
import { mediaApi } from "~/lib/api/media";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export async function loader({ request }: { request: Request }) {
  const url = new URL(request.url);
  const page = Number.parseInt(url.searchParams.get("page") || "1");
  const pageSize = Number.parseInt(url.searchParams.get("limit") || "10");
  const search = url.searchParams.get("search") || "";

  const result = await mediaApi.getAll({
    request,
    page,
    pageSize,
    search: { file_name: search },
  });

  return result;
}

export async function action({ request }: { request: Request }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const id = formData.get("id")?.toString();
  const ids = formData.getAll("ids[]");

  try {
    if (intent === "delete" && id) {
      await mediaApi.deleteById({
        request,
        id,
      });
      return { success: true, message: "Media deleted successfully" };
    }

    if (intent === "bulkDelete" && ids.length > 0) {
      await mediaApi.deleteMany({
        request,
        ids: ids.map((id) => id.toString()),
      });
      return { success: true, message: `${ids.length} media file(s) deleted successfully` };
    }

    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}

export default function MediaDashboardPage() {
  const loaderData = useLoaderData<typeof loader>();
  const actionData = useActionData<typeof action>();
  const submit = useSubmit();

  const [searchParams, setSearchParams] = useSearchParams();

  const [searchValue, setSearchValue] = useState(searchParams.get("search") || "");

  useEffect(() => {
    const success = searchParams.get("success");
    if (success === "uploaded") {
      toast.success("Media uploaded successfully!");
      const newParams = new URLSearchParams(searchParams);
      newParams.delete("success");
      setSearchParams(newParams, { replace: true });
    }
  }, [searchParams, setSearchParams]);

  useEffect(() => {
    if (actionData) {
      if (actionData.success) {
        toast.success(actionData.message || "Media deleted successfully!");
      } else {
        toast.error(actionData.message || "Failed to delete media");
      }
    }
  }, [actionData]);

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

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Media"
          description="Manage and organize your media files"
          addButtonText="Upload Media"
          addButtonLink="/dashboard/media/upload"
        />

        <DataTable
          title="All Media"
          description={`${loaderData.pagination?.totalDocs} media file${loaderData.pagination?.totalDocs !== 1 ? "s" : ""} total`}
          data={loaderData.data ?? []}
          columns={columns}
          searchPlaceholder="Search media files..."
          searchValue={searchValue}
          onSearchChange={handleSearch}
          emptyMessage="No media files found."
          totalPages={loaderData.pagination?.totalPages ?? 0}
          manualPagination
          enableRowSelection={true}
          onBulkDelete={handleBulkDelete}
          getRowId={(row) => row.id}
        />

        {(loaderData.pagination?.totalPages ?? 0) > 1 && (
          <TablePagination
            href={"/dashboard/media"}
            currentPage={loaderData.pagination?.page ?? 1}
            totalPages={loaderData.pagination?.totalPages ?? 0}
            onPageChange={handlePageChange}
          />
        )}
      </div>
    </div>
  );
}
