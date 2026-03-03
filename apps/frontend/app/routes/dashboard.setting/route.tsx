import { redirect, data } from "react-router";
import type { Route } from "./+types/route";
import { UpdateProfileForm } from "./update-profile-form";
import { ChangePasswordForm } from "./change-password-form";
import { getSession, commitSession } from "~/session.server";
import { authApi } from "~/features/users/api/auth";
import type { ActionData } from "~/types";

export function meta() {
  return [
    { title: "Settings - Starter App" },
    { name: "description", content: "Manage your account settings" },
  ];
}

export async function loader({ request }: Route.LoaderArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const userId = session.get("userId");
  const userName = session.get("userName");
  const userEmail = session.get("userEmail");

  if (!userId) {
    return redirect("/login");
  }

  return data({
    user: {
      id: userId,
      name: userName || "",
      email: userEmail || "",
    },
  });
}

export async function action({ request }: Route.ActionArgs) {
  const session = await getSession(request.headers.get("Cookie"));
  const userId = session.get("userId");
  const token = session.get("token");
  console.log("Action called in settings route : token : ", token);

  if (!userId || !token) {
    return redirect("/login");
  }

  const formData = await request.formData();
  const actionType = formData.get("actionType");

  if (actionType === "updateProfile") {
    const name = formData.get("name") as string;
    const email = formData.get("email") as string;

    if (!name || !email) {
      return data({
        success: false,
        actionType: "updateProfile",
        errors: { general: "Name and email are required" },
      } as ActionData);
    }

    try {
      const result = await authApi.updateProfile({ name, email }, token);

      if (result.code === 200 && result.data) {
        session.set("userName", result.data.name);
        session.set("userEmail", result.data.email);

        return data(
          {
            success: true,
            actionType: "updateProfile",
            errors: {},
          } as ActionData,
          {
            headers: {
              "Set-Cookie": await commitSession(session),
            },
          }
        );
      }

      return data({
        success: false,
        actionType: "updateProfile",
        errors: { general: result.message || "Failed to update profile" },
      } as ActionData);
    } catch (error) {
      return data({
        success: false,
        actionType: "updateProfile",
        errors: { general: "An error occurred while updating your profile" },
      } as ActionData);
    }
  }

  if (actionType === "changePassword") {
    const currentPassword = formData.get("currentPassword") as string;
    const newPassword = formData.get("newPassword") as string;
    const confirmPassword = formData.get("confirmPassword") as string;

    if (!currentPassword || !newPassword || !confirmPassword) {
      return data({
        success: false,
        actionType: "changePassword",
        errors: { general: "All password fields are required" },
      } as ActionData);
    }

    if (newPassword !== confirmPassword) {
      return data({
        success: false,
        actionType: "changePassword",
        errors: { general: "New passwords do not match" },
      } as ActionData);
    }

    try {
      const result = await authApi.changePassword({ currentPassword, newPassword }, token);

      if (result.code === 200) {
        return data({
          success: true,
          actionType: "changePassword",
          errors: {},
        } as ActionData);
      }

      return data({
        success: false,
        actionType: "changePassword",
        errors: { general: result.message || "Failed to change password" },
      } as ActionData);
    } catch (error) {
      return data({
        success: false,
        actionType: "changePassword",
        errors: { general: "An error occurred while changing your password" },
      } as ActionData);
    }
  }

  return data({
    success: false,
    actionType: "unknown",
    errors: { general: "Invalid action type" },
  } as ActionData);
}

export default function SettingsPage({ loaderData }: Route.ComponentProps) {
  return (
    <div className="mx-auto py-8 px-8 w-full">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
        <p className="text-muted-foreground mt-2">Manage your account settings and preferences</p>
      </div>

      <div className="space-y-6">
        <UpdateProfileForm
          defaultValues={{
            name: loaderData.user.name,
            email: loaderData.user.email,
          }}
        />
        <ChangePasswordForm />
      </div>
    </div>
  );
}
