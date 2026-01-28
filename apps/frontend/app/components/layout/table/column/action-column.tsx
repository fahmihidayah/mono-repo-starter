import type { ColumnDef } from "@tanstack/react-table";
import { Edit, MoreHorizontal, Trash2 } from "lucide-react";
import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { toast } from "sonner";

export interface ActionColumnConfig<T> {
  onEdit?: (item: T) => void;
  onDelete?: (item: T) => void;
  onCopyId?: (item: T) => void;
  getItemId: (item: T) => string;
  getItemName?: (item: T) => string;
}

export function createActionColumn<T>(config: ActionColumnConfig<T>): ColumnDef<T> {
  return {
    id: "actions",
    cell: ({ row }) => {
      const item = row.original;
      const itemId = config.getItemId(item);

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="size-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="size-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            {config.onCopyId && (
              <DropdownMenuItem
                onClick={() => {
                  navigator.clipboard.writeText(itemId);
                  toast.success("ID copied to clipboard");
                }}
              >
                Copy ID
              </DropdownMenuItem>
            )}
            {(config.onEdit || config.onDelete) && <DropdownMenuSeparator />}
            {config.onEdit && (
              <DropdownMenuItem onClick={() => config.onEdit?.(item)}>
                <Edit className="mr-2 size-4" />
                Edit
              </DropdownMenuItem>
            )}
            {config.onDelete && (
              <DropdownMenuItem
                onClick={() => config.onDelete?.(item)}
                className="text-destructive"
              >
                <Trash2 className="mr-2 size-4" />
                Delete
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  };
}
