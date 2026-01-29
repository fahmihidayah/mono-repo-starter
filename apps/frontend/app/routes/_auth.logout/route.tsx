import { redirect } from "react-router";
import { logout } from "~/lib/auth.server";
// import type { Route } from "../+types/logout";
import type { Route } from "./+types/route";
/**
 * Logout Action - Destroys session and calls API logout
 * POST /logout
 */
export async function action({ request }: Route.ActionArgs) {
  return await logout(request);
}

/**
 * Logout Loader - Redirect GET requests
 * Logout should only be done via POST for security
 */
export async function loader() {
  return redirect("/");
}
