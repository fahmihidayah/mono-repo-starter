import { useCallback, useState } from "react";
import { Upload, X } from "lucide-react";
import { cn } from "~/lib/utils";

interface ImageFieldProps {
  value?: File | null;
  onChange: (file: File | null) => void;
  accept?: string;
  maxSize?: number; // in bytes
  disabled?: boolean;
  className?: string;
  placeholder?: string;
}

export function ImageField({
  value,
  onChange,
  accept = "image/*",
  maxSize = 5 * 1024 * 1024, // 5MB default
  disabled = false,
  className,
  placeholder,
}: ImageFieldProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [preview, setPreview] = useState<string | null>(null);

  // Generate preview when value changes
  const updatePreview = useCallback((file: File | null) => {
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    } else {
      setPreview(null);
    }
  }, []);

  // Update preview when value changes
  useState(() => {
    if (value) {
      updatePreview(value);
    }
  });

  const validateFile = useCallback(
    (file: File): boolean => {
      setError(null);

      if (file.size > maxSize) {
        setError(`File must be smaller than ${(maxSize / (1024 * 1024)).toFixed(0)}MB`);
        return false;
      }

      // Validate file type
      if (accept && !file.type.match(accept.replace("*", ".*"))) {
        setError(`Invalid file type. Expected ${accept}`);
        return false;
      }

      return true;
    },
    [maxSize, accept]
  );

  const handleFileSelect = useCallback(
    (file: File | null) => {
      if (!file) {
        onChange(null);
        setPreview(null);
        setError(null);
        return;
      }

      if (validateFile(file)) {
        onChange(file);
        updatePreview(file);
      }
    },
    [validateFile, onChange, updatePreview]
  );

  const handleDragEnter = useCallback(
    (e: React.DragEvent) => {
      if (disabled) return;
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(true);
    },
    [disabled]
  );

  const handleDragLeave = useCallback(
    (e: React.DragEvent) => {
      if (disabled) return;
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(false);
    },
    [disabled]
  );

  const handleDragOver = useCallback(
    (e: React.DragEvent) => {
      if (disabled) return;
      e.preventDefault();
      e.stopPropagation();
    },
    [disabled]
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      if (disabled) return;
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(false);

      const files = e.dataTransfer.files;
      if (files.length > 0) {
        handleFileSelect(files[0]);
      }
    },
    [disabled, handleFileSelect]
  );

  const handleFileInput = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      if (disabled) return;
      const files = e.target.files;
      if (files && files.length > 0) {
        handleFileSelect(files[0]);
      }
      // Reset input
      e.target.value = "";
    },
    [disabled, handleFileSelect]
  );

  const handleRemove = useCallback(
    (e: React.MouseEvent) => {
      e.stopPropagation();
      handleFileSelect(null);
    },
    [handleFileSelect]
  );

  return (
    <div className={cn("w-full", className)}>
      <div
        onDragEnter={handleDragEnter}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        className={cn(
          "relative flex flex-col items-center justify-center w-full border-2 border-dashed rounded-lg transition-colors",
          !disabled && "cursor-pointer",
          disabled && "opacity-50 cursor-not-allowed",
          preview ? "min-h-48" : "min-h-32",
          isDragging && !disabled
            ? "border-primary bg-primary/5"
            : "border-gray-300 hover:border-primary/50 hover:bg-gray-50"
        )}
      >
        {!disabled && (
          <input
            type="file"
            accept={accept}
            onChange={handleFileInput}
            disabled={disabled}
            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
          />
        )}

        {preview ? (
          <div className="relative w-full h-full p-4">
            <img
              src={preview}
              alt="Preview"
              className="w-full h-full object-contain max-h-64 rounded"
            />
            {!disabled && (
              <button
                type="button"
                onClick={handleRemove}
                className="absolute top-2 right-2 p-1.5 bg-destructive text-destructive-foreground rounded-full hover:bg-destructive/90 transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            )}
            {value && (
              <p className="text-xs text-gray-600 mt-2 text-center truncate">{value.name}</p>
            )}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center p-6 text-center">
            <Upload
              className={cn(
                "w-10 h-10 mb-3 transition-colors",
                isDragging && !disabled ? "text-primary" : "text-gray-400"
              )}
            />
            <p className="mb-1 text-sm font-medium text-gray-700">
              {placeholder || (
                <>
                  <span className="text-primary">Click to upload</span> or drag and drop
                </>
              )}
            </p>
            <p className="text-xs text-gray-500">
              {accept === "image/*" ? "Images" : "Files"} up to {(maxSize / (1024 * 1024)).toFixed(0)}
              MB
            </p>
          </div>
        )}
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
