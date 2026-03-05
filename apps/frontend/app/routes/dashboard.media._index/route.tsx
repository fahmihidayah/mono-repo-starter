import { useLoaderData } from "react-router";
import { MediaList } from "~/features/media/components";
import { mediaListLoader } from "~/features/media/loaders/mediaListLoader";
import { deleteMediaAction } from "~/features/media/actions/deleteMedia";
import { useMediaList } from "~/features/media/hooks/useMediaList";

export async function loader({ request }: { request: Request }) {
  return mediaListLoader({ request });
}

export async function action({ request }: { request: Request }) {
  return deleteMediaAction({ request });
}

export default function MediaDashboardPage() {
  const loaderData = useLoaderData<typeof loader>();
  const { searchValue, handleSearch, handlePageChange, handleBulkDelete } = useMediaList();

  return (
    <MediaList
      media={loaderData.data ?? []}
      pagination={{
        page: loaderData.pagination?.page ?? 1,
        totalPages: loaderData.pagination?.totalPages ?? 0,
        totalDocs: loaderData.pagination?.totalDocs ?? 0,
      }}
      searchValue={searchValue}
      onSearchChange={handleSearch}
      onPageChange={handlePageChange}
      onBulkDelete={handleBulkDelete}
    />
  );
}
