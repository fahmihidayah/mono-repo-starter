import { useState } from "react";
import { DropZone } from "./DropZone";
import { ImagePreview } from "./ImagePreview";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";

interface MediaUploadFormProps {
  onSubmit: (files: File[], alt: string) => void;
  isSubmitting?: boolean;
  maxFiles?: number;
  maxSize?: number;
}

export function MediaUploadForm({
  onSubmit,
  isSubmitting = false,
  maxFiles = 10,
  maxSize = 5 * 1024 * 1024,
}: MediaUploadFormProps) {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [alt, setAlt] = useState("");

  const handleFilesSelected = (files: File[]) => {
    setSelectedFiles((prev) => {
      const combined = [...prev, ...files];
      return combined.slice(0, maxFiles);
    });
  };

  const handleRemoveFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedFiles.length > 0) {
      onSubmit(selectedFiles, alt);
    }
  };

  const handleClear = () => {
    setSelectedFiles([]);
    setAlt("");
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 relative ">
      <div className="space-y-2">
        <div className="flex justify-between items-center">
          <Label>Upload Files</Label>
          {selectedFiles.length > 0 && (
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Uploading..." : "Upload Files"}
              </Button>
              <Button type="button" variant="outline" onClick={handleClear} disabled={isSubmitting}>
                Clear All
              </Button>
            </div>
          )}
        </div>

        {selectedFiles.length === 0 && (
          <DropZone
            onFilesSelected={handleFilesSelected}
            accept="image/*"
            maxFiles={maxFiles}
            maxSize={maxSize}
          />
        )}
      </div>

      {selectedFiles.length > 0 && (
        <>
          <div className="space-y-2">
            <Label>Selected Files ({selectedFiles.length})</Label>
            <ImagePreview files={selectedFiles} onRemove={handleRemoveFile} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="alt">Alt Text (Optional)</Label>
            <Input
              id="alt"
              type="text"
              value={alt}
              onChange={(e) => setAlt(e.target.value)}
              placeholder="Describe the images for accessibility"
              disabled={isSubmitting}
            />
            <p className="text-xs text-gray-500">Provide a description for accessibility and SEO</p>
          </div>
        </>
      )}
    </form>
  );
}
