interface DetailLoaderConfig<T> {
  request: Request;
  id: string;
  getByIdFn: (id: string, request: Request) => Promise<T>;
}

export function createDetailLoader<T>({
  request,
  id,
  getByIdFn,
}: DetailLoaderConfig<T>) {
  return getByIdFn(id, request);
}
