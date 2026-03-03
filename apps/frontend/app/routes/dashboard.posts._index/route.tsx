import { useLoaderData } from "react-router";
import { postListLoader } from "~/features/posts/loaders";
import { deletePostAction } from "~/features/posts/actions";
import { PostList } from "~/features/posts/components";
import { usePostList } from "~/features/posts/hooks";

export async function loader({ request }: { request: Request }) {
  return postListLoader({ request });
}

export async function action({ request }: { request: Request }) {
  return deletePostAction({ request });
}

export default function PostsDashboardPage() {
  const loaderData = useLoaderData<typeof loader>();
  const { searchValue, handleSearch, handlePageChange, handleBulkDelete } = usePostList();

  return (
    <PostList
      posts={loaderData.data ?? []}
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
