import type { ColumnDef } from "@tanstack/react-table";
import { Edit, MoreHorizontal, Trash2 } from "lucide-react";
import { Link, Form } from "react-router";
import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";

export interface ActionColumnConfig<T> {
  getItemId: (item: T) => string;
  editLink?: (item: T) => string;
  deleteFormAction?: string;
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
            <DropdownMenuSeparator />
            {config.editLink && (
              <DropdownMenuItem asChild>
                <Link to={config.editLink(item)}>
                  <Edit className="mr-2 size-4" />
                  Edit
                </Link>
              </DropdownMenuItem>
            )}
            {config.deleteFormAction && (
              <Form method="post" action={config.deleteFormAction}>
                <input type="hidden" name="intent" value="delete" />
                <input type="hidden" name="id" value={itemId} />
                <DropdownMenuItem asChild>
                  <button type="submit" className="w-full text-destructive">
                    <Trash2 className="mr-2 size-4" />
                    Delete
                  </button>
                </DropdownMenuItem>
              </Form>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      );
    },
  };
}
