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
      <div>
        <h1 className="text-3xl font-bold">Upload Media</h1>
      </div>
      <div className="py-5">
        <MediaUploadForm onSubmit={handleUpload} isSubmitting={isSubmitting} />
      </div>
    </div>
  );
}
