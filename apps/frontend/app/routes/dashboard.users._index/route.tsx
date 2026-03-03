import { useLoaderData } from "react-router";
import { userListLoader } from "~/features/users/loaders";
import { deleteUserAction } from "~/features/users/actions";
import { UserList } from "~/features/users/components";
import { useUserList } from "~/features/users/hooks";

export async function loader({ request }: { request: Request }) {
  return userListLoader({ request });
}

export async function action({ request }: { request: Request }) {
  return deleteUserAction({ request });
}

export function meta() {
  return [{ title: "User - Dashboard" }, { name: "description", content: "Manage your tasks" }];
}

export default function DashboardUsersPage() {
  const loaderData = useLoaderData<typeof loader>();
  const { searchValue, handleSearch, handlePageChange, handleBulkDelete } = useUserList();

  return (
    <UserList
      users={loaderData.data ?? []}
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
