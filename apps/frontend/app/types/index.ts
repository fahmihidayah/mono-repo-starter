

// Define the action data type
export type ActionData = {
  success: boolean;
  errors: { [key: string]: string | undefined };
  actionType?: string;
};

export type BaseResponse<T> = {
  message: string,
  code: number,
  data?: T,
  pagination?: {
    limit: number,
    nextPage: number | null,
    page: number,
    prevPage: number | null,
    totalDocs: number,
    totalPages: number
  }
}
