import { useResourceList } from "~/lib/hooks";

export function useMediaList() {
  return useResourceList({
    successMessage: "Media uploaded successfully!",
  });
}
