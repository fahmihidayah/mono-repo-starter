import { useLoaderData } from "react-router";
import { categoryListLoader } from "~/features/category/loaders";
import { deleteCategoryAction } from "~/features/category/actions";
import { CategoryList } from "~/features/category/components";
import { useCategoryList } from "~/features/category/hooks";

export async function loader({ request }: { request: Request }) {
  return categoryListLoader({ request });
}

export async function action({ request }: { request: Request }) {
  return deleteCategoryAction({ request });
}

export default function CategoriesDashboardPage() {
  const loaderData = useLoaderData<typeof loader>();
  const { searchValue, handleSearch, handlePageChange, handleBulkDelete } = useCategoryList();

  return (
    <CategoryList
      categories={loaderData.data ?? []}
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
