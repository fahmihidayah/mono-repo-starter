import type { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router";

export interface UrlColumnConfig<T> {
  type: "url";
  accessorKey: keyof T;
  header: string;
  className?: string;
  fallback?: string;
  format?: (value: any) => string;
  maxLength?: number;
}

export function createUrlColumn<T>(config: UrlColumnConfig<T>): ColumnDef<T> {
  return {
    accessorKey: config.accessorKey as string,
    header: config.header,
    cell: ({ row }) => {
      const value = row.getValue(config.accessorKey as string);
      const url = config.format ? config.format(value) : value;

      if (!url || typeof url !== "string") {
        return <span className="text-muted-foreground text-sm">{config.fallback || "No URL"}</span>;
      }

      const displayText =
        config.maxLength && url.length > config.maxLength
          ? `${url.substring(0, config.maxLength)}...`
          : url;

      const classNames = [
        "text-blue-600 hover:text-blue-800 underline text-sm truncate block",
        config.className,
      ].filter(Boolean);

      const className = classNames.length > 0 ? classNames.join(" ") : undefined;

      const isExternal = url.startsWith("http://") || url.startsWith("https://");

      if (isExternal) {
        return (
          <a href={url} target="_blank" rel="noopener noreferrer" className={className} title={url}>
            {displayText}
          </a>
        );
      }

      return (
        <Link to={url} className={className} title={url}>
          {displayText}
        </Link>
      );
    },
  };
}
