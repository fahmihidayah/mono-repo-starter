import { MediaForm } from "~/features/media/components";
import { createMediaAction } from "~/features/media/actions/createMedia";
import { useResourceForm } from "~/lib/hooks";

export async function action({ request }: { request: Request }) {
  return createMediaAction({ request });
}

export default function MediaAddPage() {
  const { submitFormSchema, isSubmitting } = useResourceForm();

  return (
    <div className="container w-full mx-auto p-5">
      <div>
        <h1 className="text-3xl font-bold">Upload Media</h1>
      </div>
      <div className="py-5">
        <MediaForm
          defaultValues={{ alt: "" }}
          onSubmit={submitFormSchema}
          isSubmitting={isSubmitting}
          submitButtonText="Upload"
        />
      </div>
    </div>
  );
}
