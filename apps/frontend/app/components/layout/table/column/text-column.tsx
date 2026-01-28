import type { ColumnDef } from "@tanstack/react-table";

export interface TextColumnConfig<T> {
  type : "text"
  accessorKey: keyof T;
  header: string;
  className?: string;
  fallback?: string;
  format?: (value: any) => string;
  isBold?: boolean;
}

export function createTextColumn<T>(config: TextColumnConfig<T>): ColumnDef<T> {
  return {
    accessorKey: config.accessorKey as string,
    header: config.header,
    cell: ({ row }) => {
      const value = row.getValue(config.accessorKey as string);
      const displayValue = config.format
        ? config.format(value)
        : value || config.fallback || "-";

      const classNames = [
        config.className,
        config.isBold && "font-medium",
      ].filter(Boolean);

      const className = classNames.length > 0 ? classNames.join(" ") : undefined;

      return <div className={className}>{String(displayValue)}</div>;
    },
  };
}
