import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/components/ui/card";
import { authClient } from "~/lib/auth-client";
import { toast } from "sonner";
export function meta() {
  return [
    { title: "Register - Starter App" },
    { name: "description", content: "Create a new account" },
  ];
}

export default function Register() {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    console.log("[Register] Form submitted", { name: formData.name, email: formData.email });

    // Validate passwords match
    if (formData.password !== formData.confirmPassword) {
      toast.error("Registration Failed", {
        description: "Passwords do not match",
      });
      return;
    }

    // Validate password length
    if (formData.password.length < 8) {
      toast.error("Registration Failed", {
        description: "Password must be at least 8 characters",
      });
      return;
    }

    setIsLoading(true);

    try {
      console.log("[Register] Attempting sign up...");
      const response = await authClient.signUp.email({
        name: formData.name,
        email: formData.email,
        password: formData.password,
      });

      console.log("[Register] Sign up response:", response);

      if (response.error) {
        console.error("[Register] Sign up failed:", response.error);
        toast.error("Registration Failed", {
          description: response.error.message || "Unable to create account",
        });
        return;
      }

      console.log("[Register] Sign up successful, redirecting to dashboard");
      toast.success("Registration Successful", {
        description: "Welcome! Your account has been created.",
      });

      navigate("/dashboard");
    } catch (err: any) {
      console.error("[Register] Unexpected error:", err);
      toast.error("Registration Failed", {
        description: err.message || "An unexpected error occurred",
      });
    } finally {
      setIsLoading(false);
      console.log("[Register] Loading state set to false");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-background via-background to-muted/20 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">Create an account</CardTitle>
          <CardDescription className="text-center">
            Enter your information to get started
          </CardDescription>
        </CardHeader>
        <form>
          <CardContent className="flex flex-col gap-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                type="text"
                placeholder="John Doe"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="john.doe@example.com"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                required
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                required
                disabled={isLoading}
                minLength={8}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Confirm Password</Label>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="••••••••"
                value={formData.confirmPassword}
                onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
                required
                disabled={isLoading}
                minLength={8}
              />
            </div>
          </CardContent>
          <CardFooter className="flex flex-col pt-5 space-y-5">
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Creating account..." : "Create account"}
            </Button>
            <div className="text-sm text-center text-muted-foreground">
              Already have an account?{" "}
              <Link to="/login" className="text-primary hover:underline font-medium">
                Sign in
              </Link>
            </div>
            <div className="text-sm text-center text-muted-foreground">
              <Link to="/" className="text-primary hover:underline font-medium">
                Back to home
              </Link>
            </div>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
