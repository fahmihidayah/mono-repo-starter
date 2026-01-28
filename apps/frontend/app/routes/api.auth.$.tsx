// Stub auth route - replace with actual implementation
import type { Route } from "./+types/api.auth.$";

export async function loader({ params }: Route.LoaderArgs) {
  // TODO: Implement auth handler with Go API
  return new Response("Auth not implemented", { status: 501 });
}

export async function action({ request, params }: Route.ActionArgs) {
  // TODO: Implement auth handler with Go API
  return new Response("Auth not implemented", { status: 501 });
}
