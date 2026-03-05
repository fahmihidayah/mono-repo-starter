import { MediaUploadForm } from "~/features/media/components";
import { uploadMediaAction } from "~/features/media/actions/uploadMedia";
import { useMediaUpload } from "~/features/media/hooks/useMediaUpload";

export async function action({ request }: { request: Request }) {
  return uploadMediaAction({ request });
}

export default function MediaAddPage() {
  const { handleUpload, isSubmitting } = useMediaUpload();

  return (
    <div className="flex-1 p-6">
      <div className="max-w-4xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Upload Media</h1>
          <p className="text-muted-foreground mt-2">
            Upload images and other media files to your library
          </p>
        </div>

        <div className="bg-white rounded-lg border p-6">
          <MediaUploadForm onSubmit={handleUpload} isSubmitting={isSubmitting} />
        </div>
      </div>
    </div>
  );
}
