import { redirect } from "react-router";
import type { Route } from "./+types/route";
import { authApi } from "~/features/users/api/auth";
import * as serverSession from "~/session.server";

/**
 * Logout Action - Destroys session and calls API logout
 * POST /logout
 */
export async function action({ request }: Route.ActionArgs) {
  const { getSession, commitSession, destroySession } = serverSession;
  const session = await getSession(request.headers.get("Cookie"));
  const token = session.get("token");

  if (!token) {
    return redirect("/");
  }

  try {
    const result = await authApi.logout(token);

    if (result.code !== 200) {
      // Set flash error message
      session.flash("error", "Failed to logout: " + (result.message || "Unknown error"));

      return redirect("/dashboard", {
        headers: {
          "Set-Cookie": await commitSession(session),
        },
      });
    }

    console.log("Logout action successful");
    // Logout successful - destroy session
    return redirect("/", {
      headers: {
        "Set-Cookie": await destroySession(session),
      },
    });
  } catch (error: any) {
    console.log("Logout action error: ", error);
    // Set flash error message for network/API errors
    session.flash("error", "Failed to logout: " + (error.message || "Network error"));

    return redirect("/dashboard", {
      headers: {
        "Set-Cookie": await commitSession(session),
      },
    });
  }
}

/**
 * Logout Loader - Redirect GET requests
 * Logout should only be done via POST for security
 */
export async function loader() {
  return redirect("/");
}
