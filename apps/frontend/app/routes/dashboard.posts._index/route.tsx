import type { Post } from "~/features/posts/type";
import createColumn from "~/components/layout/table/column/create-column";
import { DataTable, PageHeader, TablePagination } from "~/components/layout/table/table-list";
import { useLoaderData, useSearchParams, useActionData, useSubmit } from "react-router";
import { postApi } from "~/lib/api/posts";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export async function loader({ request }: { request: Request }) {
  const url = new URL(request.url);
  const page = Number.parseInt(url.searchParams.get("page") || "1");
  const pageSize = Number.parseInt(url.searchParams.get("limit") || "10");
  const search = url.searchParams.get("search") || "";

  const result = await postApi.getAll({
    request,
    page,
    pageSize,
    search: { title: search },
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
      await postApi.deleteById({
        request,
        id,
      });
      return { success: true, message: "Post deleted successfully" };
    }

    if (intent === "bulkDelete" && ids.length > 0) {
      await postApi.deleteMany({
        request,
        ids: ids.map((id) => id.toString()),
      });
      return { success: true, message: `${ids.length} post${ids.length !== 1 ? "s" : ""} deleted successfully` };
    }

    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}

export default function PostsDashboardPage() {
  const loaderData = useLoaderData<typeof loader>();
  const actionData = useActionData<typeof action>();
  const submit = useSubmit();

  const [searchParams, setSearchParams] = useSearchParams();

  const [searchValue, setSearchValue] = useState(searchParams.get("search") || "");

  useEffect(() => {
    const success = searchParams.get("success");
    if (success === "created") {
      toast.success("Post created successfully!");
      const newParams = new URLSearchParams(searchParams);
      newParams.delete("success");
      setSearchParams(newParams, { replace: true });
    }
  }, [searchParams, setSearchParams]);

  useEffect(() => {
    if (actionData) {
      if (actionData.success) {
        toast.success(actionData.message || "Post deleted successfully!");
      } else {
        toast.error(actionData.message || "Failed to delete post");
      }
    }
  }, [actionData]);

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
          title="Posts"
          description="Manage and organize your posts efficiently"
          addButtonText="Add Post"
          addButtonLink="/dashboard/posts/add"
        />

        <DataTable
          title="All Posts"
          description={`${loaderData.pagination?.totalDocs} post${loaderData.pagination?.totalDocs !== 1 ? "s" : ""} total`}
          data={loaderData.data ?? []}
          columns={columns}
          enableRowSelection={true}
          onBulkDelete={handleBulkDelete}
          getRowId={(row) => row.id || ""}
          searchPlaceholder="Search posts..."
          searchValue={searchValue}
          onSearchChange={handleSearch}
          emptyMessage="No posts found."
          totalPages={loaderData.pagination?.totalPages ?? 0}
          manualPagination
        />

        {(loaderData.pagination?.totalPages ?? 0) > 1 && (
          <TablePagination
            href={"/dashboard/posts"}
            currentPage={loaderData.pagination?.page ?? 1}
            totalPages={loaderData.pagination?.totalPages ?? 0}
            onPageChange={handlePageChange}
          />
        )}
      </div>
    </div>
  );
}
