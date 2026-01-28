# Table List Components

Reusable components for building table-based list views with headers, search, pagination, and dialogs.

## Structure

```
table-list/
├── index.tsx              # Main exports
├── page-header.tsx        # Page title and action button
├── data-table.tsx         # Reusable data table with search
├── table-pagination.tsx   # Pagination controls
├── delete-dialog.tsx      # Reusable delete confirmation dialog
└── README.md             # This file
```

## Components

### PageHeader

Displays a page title, description, and an optional "Add" button.

**Props:**
- `title: string` - Page title
- `description?: string` - Optional description text
- `addButtonText?: string` - Button text (default: "Add New")
- `addButtonLink?: string` - Link URL for the button
- `onAddClick?: () => void` - Click handler (if no link)

**Usage:**
```tsx
<PageHeader
  title="Tasks"
  description="Manage your tasks"
  addButtonText="Add Task"
  addButtonLink="/dashboard/tasks/add"
/>
```

### DataTable

A reusable table component with built-in search functionality.

**Props:**
- `title: string` - Table card title
- `description?: string` - Table card description
- `data: TData[]` - Array of data items
- `columns: ColumnDef<TData>[]` - TanStack Table column definitions
- `searchable?: boolean` - Enable search (default: true)
- `searchPlaceholder?: string` - Search input placeholder
- `searchValue?: string` - Controlled search value
- `onSearchChange?: (value: string) => void` - Search change handler
- `emptyMessage?: string` - Message when no data
- `totalPages?: number` - For manual pagination
- `manualPagination?: boolean` - Use manual pagination (default: false)

**Usage:**
```tsx
<DataTable
  title="All Tasks"
  description="10 tasks total"
  data={tasks}
  columns={columns}
  searchPlaceholder="Search tasks..."
  searchValue={searchValue}
  onSearchChange={handleSearch}
  totalPages={totalPages}
  manualPagination
/>
```

### TablePagination

Pagination controls with previous/next and page numbers.

**Props:**
- `currentPage: number` - Current page number (1-indexed)
- `totalPages: number` - Total number of pages
- `onPageChange: (page: number) => void` - Page change handler

**Usage:**
```tsx
<TablePagination
  currentPage={page}
  totalPages={totalPages}
  onPageChange={handlePageChange}
/>
```

### DeleteDialog

A reusable confirmation dialog for deletions.

**Props:**
- `open: boolean` - Dialog open state
- `onOpenChange: (open: boolean) => void` - Open state change handler
- `title?: string` - Dialog title (default: "Delete Item")
- `description?: string` - Custom description
- `itemName?: string` - Name of item to delete (auto-generates description)
- `onConfirm: () => void` - Confirm button handler
- `onCancel?: () => void` - Cancel button handler
- `confirmButtonText?: string` - Confirm button text (default: "Delete")
- `cancelButtonText?: string` - Cancel button text (default: "Cancel")

**Usage:**
```tsx
<DeleteDialog
  open={!!deletingItem}
  onOpenChange={(open) => !open && setDeletingItem(null)}
  title="Delete Task"
  itemName={deletingItem?.title}
  onConfirm={handleDelete}
/>
```

## Complete Example

Here's how all components work together (see [dashboard.tasks.tsx](../../../routes/dashboard.tasks.tsx)):

```tsx
import { PageHeader, DataTable, TablePagination, DeleteDialog } from "~/components/layout/table-list";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "~/components/ui/dialog";

export default function TasksPage() {
  const { tasks, page, totalPages } = useLoaderData<typeof loader>();
  const [searchValue, setSearchValue] = useState("");
  const [deletingTask, setDeletingTask] = useState<Task | null>(null);

  const columns: ColumnDef<Task>[] = [
    {
      accessorKey: "title",
      header: "Title",
      cell: ({ row }) => <div>{row.getValue("title")}</div>,
    },
    // ... more columns
  ];

  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Tasks"
          description="Manage your tasks"
          addButtonText="Add Task"
          addButtonLink="/tasks/add"
        />

        <DataTable
          title="All Tasks"
          data={tasks}
          columns={columns}
          searchValue={searchValue}
          onSearchChange={setSearchValue}
          totalPages={totalPages}
          manualPagination
        />

        <TablePagination
          currentPage={page}
          totalPages={totalPages}
          onPageChange={handlePageChange}
        />
      </div>

      {/* Use standard Dialog component for edit dialogs */}
      <Dialog open={!!editingTask} onOpenChange={() => setEditingTask(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Task</DialogTitle>
          </DialogHeader>
          {/* Your custom edit form */}
        </DialogContent>
      </Dialog>

      <DeleteDialog
        open={!!deletingTask}
        onOpenChange={(open) => !open && setDeletingTask(null)}
        itemName={deletingTask?.title}
        onConfirm={handleDelete}
      />
    </div>
  );
}
```

## Benefits

1. **Reusable** - Use across multiple list pages
2. **Consistent UI** - Same look and feel everywhere
3. **Less Code** - Reduced boilerplate
4. **Type Safe** - Full TypeScript support
5. **Customizable** - Props for common customizations
6. **TanStack Table** - Built-in table features (sorting, filtering, pagination)

## Customization

### Custom Columns

Define columns with TanStack Table's `ColumnDef`:

```tsx
const columns: ColumnDef<Task>[] = [
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const status = row.getValue("status");
      return <Badge variant={status === "done" ? "success" : "warning"}>{status}</Badge>;
    },
  },
  {
    id: "actions",
    cell: ({ row }) => (
      <DropdownMenu>
        {/* Custom actions */}
      </DropdownMenu>
    ),
  },
];
```

### No Search

Disable search in the DataTable:

```tsx
<DataTable
  data={data}
  columns={columns}
  searchable={false}
/>
```

### Custom Empty Message

```tsx
<DataTable
  data={data}
  columns={columns}
  emptyMessage="No tasks yet. Create your first task!"
/>
```

## Future Enhancements

Potential additions:
- Bulk actions (checkboxes, select all)
- Column visibility toggle
- Export to CSV/Excel
- Advanced filters
- Sortable columns
- Row expansion
