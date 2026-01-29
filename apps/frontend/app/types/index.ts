

// Define the action data type
export type ActionData = {
  success: boolean;
  errors: { [key: string]: string | undefined };
  actionType?: string;
};
