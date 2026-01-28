import { redirect } from "react-router";
import { getSession, destroySession } from "~/session.server";
import { callLogoutApi } from "~/lib/auth.server";
import type { Route } from "./+types/logout";

/**
 * Logout Action - Destroys session and calls API logout
 * POST /logout
 */
export async function action({ request }: Route.ActionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const token = session.get("token");

  // Call API logout endpoint
  if (token) {
    await callLogoutApi(token);
  }

  // Destroy session and redirect to login
  return redirect("/login", {
    headers: {
      "Set-Cookie": await destroySession(session),
    },
  });
}

/**
 * Logout Loader - Redirect GET requests
 * Logout should only be done via POST for security
 */
export async function loader() {
  return redirect("/");
}
