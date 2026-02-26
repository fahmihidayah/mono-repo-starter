import type { ColumnDef } from "@tanstack/react-table";
import { createDateColumn, type DateColumnConfig } from "./date-column";
import { createTextColumn, type TextColumnConfig } from "./text-column";
import { createActionColumn, type ActionColumnConfig } from "./action-column";
import { createImageColumn, type ImageColumnConfig } from "./image-column";
import { createUrlColumn, type UrlColumnConfig } from "./url-column";

export type CustomColumnConfig<T> = ColumnDef<T> & {
  type: "custom";
};

interface CreateColumnProps<T> {
  columnConfig: (
    | string
    | TextColumnConfig<T>
    | DateColumnConfig<T>
    | ImageColumnConfig<T>
    | UrlColumnConfig<T>
    | CustomColumnConfig<T>
  )[];
  actionColumnConfig?: ActionColumnConfig<T>;
}

export default function createColumn<T>(props: CreateColumnProps<T>): ColumnDef<T>[] {
  const columns: ColumnDef<T>[] = [];

  props.columnConfig.forEach((e) => {
    if (typeof e === "string") {
      // Capitalize first letter for header
      const header = e.charAt(0).toUpperCase() + e.slice(1);
      columns.push(
        createTextColumn<T>({
          accessorKey: e as any,
          header: header,
          type: "text",
        })
      );
    } else if (e.type === "text") {
      columns.push(createTextColumn(e));
    } else if (e.type === "date") {
      columns.push(createDateColumn(e));
    } else if (e.type === "image") {
      columns.push(createImageColumn(e));
    } else if (e.type === "url") {
      columns.push(createUrlColumn(e));
    } else if (e.type === "custom") {
      const { type, ...columnDef } = e;
      columns.push(columnDef as ColumnDef<T>);
    }
  });

  return [
    ...columns,
    props.actionColumnConfig ? createActionColumn(props.actionColumnConfig) : null,
  ].filter(Boolean) as ColumnDef<T>[];
}
