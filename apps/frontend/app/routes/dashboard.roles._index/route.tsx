import { useLoaderData } from "react-router";
import { roleListLoader } from "~/features/roles/loaders";
import { deleteRoleAction } from "~/features/roles/actions";
import { RoleList } from "~/features/roles/components";
import { useResourceList } from "~/lib/hooks";

export async function loader({ request }: { request: Request }) {
  return roleListLoader({ request });
}

export async function action({ request }: { request: Request }) {
  return deleteRoleAction({ request });
}

export default function RolesDashboardPage() {
  const loaderData = useLoaderData<typeof loader>();
  const { searchValue, handleSearch, handlePageChange, handleBulkDelete } = useResourceList({
    successMessage: "Role created successfully!",
  });

  return (
    <RoleList
      roles={loaderData.data ?? []}
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
