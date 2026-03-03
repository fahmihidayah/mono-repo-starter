interface ListLoaderConfig<T> {
  request: Request;
  getAllFn: (params: {
    request: Request;
    page: number;
    pageSize: number;
    search: Record<string, string>;
  }) => Promise<T>;
  searchField: string;
}

export function createListLoader<T>({
  request,
  getAllFn,
  searchField,
}: ListLoaderConfig<T>) {
  const url = new URL(request.url);
  const page = Number.parseInt(url.searchParams.get("page") || "1");
  const pageSize = Number.parseInt(url.searchParams.get("limit") || "10");
  const search = url.searchParams.get("search") || "";

  return getAllFn({
    request,
    page,
    pageSize,
    search: { [searchField]: search },
  });
}
