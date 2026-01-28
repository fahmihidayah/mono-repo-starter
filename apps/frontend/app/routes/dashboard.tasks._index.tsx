import { useState } from "react";
import { useLoaderData, useNavigate, useSearchParams, useSubmit } from "react-router";
import { toast } from "sonner";
import { taskRepository } from "~/features/tasks/task-repository";
import { PageHeader, DataTable, TablePagination, DeleteDialog } from "~/components/layout/table/table-list";
import { createActionColumn, createDateColumn, createTextColumn } from "~/components/layout/table/column";
import type { Route } from "./+types/dashboard.tasks";
import type { Task } from "~/features/tasks/type";
import createColumn from "~/components/layout/table/column/create-column";

// Loader - Fetch tasks with pagination and search using TaskRepository
export async function loader({ request }: Route.LoaderArgs) {
  const url = new URL(request.url);
  const page = Number.parseInt(url.searchParams.get("page") || "1");
  const pageSize = Number.parseInt(url.searchParams.get("pageSize") || "10");
  const search = url.searchParams.get("search") || "";

  // Use repository's findManyPaginated method
  const result = await taskRepository.findManyPaginated({
    page,
    pageSize,
  });

  return {
    tasks: result.data,
    totalCount: result.pagination.totalItems,
    page: result.pagination.currentPage,
    pageSize: result.pagination.pageSize,
    totalPages: result.pagination.totalPages,
  };
}

// Action - Handle delete and update operations using TaskRepository
export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const taskId = formData.get("taskId")?.toString();

  try {
    if (intent === "delete" && taskId) {
      await taskRepository.delete(taskId);
      return { success: true, message: "Task deleted successfully" };
    }

    if (intent === "update" && taskId) {
      const title = formData.get("title")?.toString();
      const description = formData.get("description")?.toString();

      await taskRepository.update(taskId, {
        title,
        description,
        updatedAt: new Date(),
      });

      return { success: true, message: "Task updated successfully" };
    }

    return { success: false, message: "Invalid action" };
  } catch (error) {
    console.error("Action error:", error);
    return { success: false, message: "An error occurred" };
  }
}

export function meta() {
  return [
    { title: "Tasks - Dashboard" },
    { name: "description", content: "Manage your tasks" },
  ];
}

export default function DashboardTasksPage() {
  const loaderData = useLoaderData<typeof loader>();
  const [searchParams, setSearchParams] = useSearchParams();

  const submit = useSubmit();

  // State
  const [searchValue, setSearchValue] = useState(searchParams.get("search") || "");
  const [deletingTask, setDeletingTask] = useState<Task | null>(null);
  const navigate = useNavigate()

  // Table columns
  const columns = createColumn<Task>({
    columnConfig: [
      // {
      //   type: "text",
      //   accessorKey: "title",
      //   header: "Title",
      //   fallback: "Untitled",
      //   isBold: true,
      // },
      "title",
      "description",
      // {
      //   type: "text",
      //   accessorKey: "description",
      //   header: "Description",
      //   className: "max-w-[500px] truncate text-muted-foreground",
      //   fallback: "No description",
      // },
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
      getItemId: (task) => task.id,
      onCopyId: () => { },
      onEdit: (task) => navigate(`${task.id}`),
      onDelete: (task) => setDeletingTask(task),
    }
  })

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

  // Handle delete
  const handleDelete = () => {
    if (!deletingTask) return;

    const formData = new FormData();
    formData.append("intent", "delete");
    formData.append("taskId", deletingTask.id);

    submit(formData, { method: "post" });
    setDeletingTask(null);
    toast.success("Task deleted successfully");
  };

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        {/* Page Header */}
        <PageHeader
          title="Tasks"
          description="Manage and organize your tasks efficiently"
          addButtonText="Add Task"
          addButtonLink="/dashboard/tasks/add"
        />

        {/* Data Table */}
        <DataTable
          title="All Tasks"
          description={`${loaderData.totalCount} task${loaderData.totalCount !== 1 ? "s" : ""} total`}
          data={loaderData.tasks}
          columns={columns}
          searchPlaceholder="Search tasks..."
          searchValue={searchValue}
          onSearchChange={handleSearch}
          emptyMessage="No tasks found."
          totalPages={loaderData.totalPages}
          manualPagination
        />

        {/* Table Pagination */}
        {loaderData.totalPages > 1 && (
          <TablePagination
            currentPage={loaderData.page}
            totalPages={loaderData.totalPages}
            onPageChange={handlePageChange}
          />
        )}
      </div>
      <DeleteDialog
        open={!!deletingTask}
        onOpenChange={(open) => !open && setDeletingTask(null)}
        title="Delete Task"
        itemName={deletingTask?.title || ""}
        onConfirm={handleDelete}
        onCancel={() => setDeletingTask(null)}
      />
    </div>
  );
}
