import { useCallback, useState } from "react";
import { Upload, X } from "lucide-react";
import { cn } from "~/lib/utils";

interface DropZoneProps {
  onFilesSelected: (files: File[]) => void;
  accept?: string;
  maxFiles?: number;
  maxSize?: number; // in bytes
  className?: string;
}

export function DropZone({
  onFilesSelected,
  accept = "image/*",
  maxFiles = 10,
  maxSize = 5 * 1024 * 1024, // 5MB default
  className,
}: DropZoneProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const validateFiles = useCallback(
    (files: FileList | File[]): File[] => {
      const fileArray = Array.from(files);
      setError(null);

      if (fileArray.length > maxFiles) {
        setError(`Maximum ${maxFiles} files allowed`);
        return [];
      }

      const oversizedFiles = fileArray.filter((file) => file.size > maxSize);
      if (oversizedFiles.length > 0) {
        setError(`Files must be smaller than ${(maxSize / (1024 * 1024)).toFixed(0)}MB`);
        return [];
      }

      return fileArray;
    },
    [maxFiles, maxSize]
  );

  const handleDragEnter = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }, []);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(false);

      const files = e.dataTransfer.files;
      const validFiles = validateFiles(files);

      if (validFiles.length > 0) {
        onFilesSelected(validFiles);
      }
    },
    [validateFiles, onFilesSelected]
  );

  const handleFileInput = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const files = e.target.files;
      if (!files) return;

      const validFiles = validateFiles(files);
      if (validFiles.length > 0) {
        onFilesSelected(validFiles);
      }

      // Reset input
      e.target.value = "";
    },
    [validateFiles, onFilesSelected]
  );

  return (
    <div className={cn("w-full", className)}>
      <div
        onDragEnter={handleDragEnter}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        className={cn(
          "relative flex flex-col items-center justify-center w-full min-h-64 border-2 border-dashed rounded-lg transition-colors cursor-pointer",
          isDragging
            ? "border-primary bg-primary/5"
            : "border-gray-300 hover:border-primary/50 hover:bg-gray-50"
        )}
      >
        <input
          type="file"
          accept={accept}
          multiple={maxFiles > 1}
          onChange={handleFileInput}
          className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        />

        <div className="flex flex-col items-center justify-center p-6 text-center">
          <Upload
            className={cn(
              "w-12 h-12 mb-4 transition-colors",
              isDragging ? "text-primary" : "text-gray-400"
            )}
          />
          <p className="mb-2 text-sm font-medium text-gray-700">
            <span className="text-primary">Click to upload</span> or drag and drop
          </p>
          <p className="text-xs text-gray-500">
            {accept === "image/*" ? "Images" : "Files"} up to{" "}
            {(maxSize / (1024 * 1024)).toFixed(0)}MB
          </p>
          {maxFiles > 1 && (
            <p className="text-xs text-gray-500 mt-1">Maximum {maxFiles} files</p>
          )}
        </div>
      </div>

      {error && (
        <div className="mt-2 flex items-center gap-2 text-sm text-destructive">
          <X className="w-4 h-4" />
          <span>{error}</span>
        </div>
      )}
    </div>
  );
}
