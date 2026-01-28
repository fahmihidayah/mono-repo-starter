import type { ColumnDef } from "@tanstack/react-table";

export interface DateColumnConfig<T> {
  type : "date"
  accessorKey: keyof T;
  header: string;
  formatDate?: (date: Date) => string;
}

export function createDateColumn<T>(config: DateColumnConfig<T>): ColumnDef<T> {
  const defaultFormatDate = (date: Date) => date.toLocaleDateString();
  const formatDate = config.formatDate || defaultFormatDate;

  return {
    accessorKey: config.accessorKey as string,
    header: config.header,
    cell: ({ row }) => {
      const date = row.getValue(config.accessorKey as string) as Date | null;
      return date ? formatDate(new Date(date)) : "-";
    },
  };
}
