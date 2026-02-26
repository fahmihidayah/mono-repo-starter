import type { ColumnDef } from "@tanstack/react-table";

export interface ImageColumnConfig<T> {
  type: "image";
  accessorKey: keyof T;
  header: string;
  className?: string;
  fallback?: string;
  format?: (value: any) => string;
  isBold?: boolean;
}

export function createImageColumn<T>(config: ImageColumnConfig<T>): ColumnDef<T> {
  return {
    accessorKey: config.accessorKey as string,
    header: config.header,
    cell: ({ row }) => {
      const value = row.getValue(config.accessorKey as string);
      const imageUrl = config.format ? config.format(value) : value;

      if (!imageUrl || typeof imageUrl !== "string") {
        return (
          <div className="flex items-center justify-center w-16 h-16 bg-muted rounded text-xs text-muted-foreground">
            No image
          </div>
        );
      }

      const imageClassNames = [
        "max-w-full max-h-full object-contain rounded",
        config.className,
      ].filter(Boolean);

      const imageClassName = imageClassNames.length > 0 ? imageClassNames.join(" ") : undefined;

      return (
        <div className="flex items-center justify-center w-16 h-16">
          <img
            src={imageUrl}
            alt={config.header}
            className={imageClassName}
            onError={(e) => {
              const img = e.target as HTMLImageElement;
              if (config.fallback) {
                img.src = config.fallback;
              } else {
                img.style.display = "none";
                const parent = img.parentElement;
                if (parent) {
                  parent.innerHTML = `<div class="flex items-center justify-center w-full h-full bg-muted rounded text-xs text-muted-foreground">No image</div>`;
                }
              }
            }}
          />
        </div>
      );
    },
  };
}
