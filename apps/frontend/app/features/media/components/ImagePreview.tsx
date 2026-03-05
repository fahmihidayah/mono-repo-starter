import { X } from "lucide-react";
import { cn } from "~/lib/utils";

interface ImagePreviewProps {
  files: File[];
  onRemove: (index: number) => void;
  className?: string;
}

export function ImagePreview({ files, onRemove, className }: ImagePreviewProps) {
  if (files.length === 0) return null;

  return (
    <div className={cn("w-full", className)}>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {files.map((file, index) => (
          <FilePreviewCard
            key={`${file.name}-${index}`}
            file={file}
            onRemove={() => onRemove(index)}
          />
        ))}
      </div>
    </div>
  );
}

interface FilePreviewCardProps {
  file: File;
  onRemove: () => void;
}

function FilePreviewCard({ file, onRemove }: FilePreviewCardProps) {
  const isImage = file.type.startsWith("image/");
  const fileSize = formatFileSize(file.size);

  return (
    <div className="relative group rounded-lg border border-gray-200 overflow-hidden bg-white hover:shadow-md transition-shadow">
      <div className="aspect-square relative bg-gray-100 flex items-center justify-center">
        {isImage ? (
          <img
            src={URL.createObjectURL(file)}
            alt={file.name}
            className="w-full h-full object-cover"
            onLoad={(e) => {
              URL.revokeObjectURL((e.target as HTMLImageElement).src);
            }}
          />
        ) : (
          <div className="flex flex-col items-center justify-center p-4 text-center">
            <div className="w-12 h-12 mb-2 rounded-lg bg-gray-200 flex items-center justify-center">
              <span className="text-xs font-medium text-gray-600">
                {file.type.split("/")[1]?.toUpperCase() || "FILE"}
              </span>
            </div>
          </div>
        )}

        <button
          type="button"
          onClick={onRemove}
          className="absolute top-2 right-2 p-1 bg-destructive text-destructive-foreground rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-destructive/90"
          aria-label={`Remove ${file.name}`}
        >
          <X className="w-4 h-4" />
        </button>
      </div>

      <div className="p-2">
        <p className="text-xs font-medium text-gray-700 truncate" title={file.name}>
          {file.name}
        </p>
        <p className="text-xs text-gray-500">{fileSize}</p>
      </div>
    </div>
  );
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${Number.parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`;
}
