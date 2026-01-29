import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function calculateSessionMaxAge(exp?: number): number {
  if (!exp) {
    // Default to 7 days
    return 60 * 60 * 24 * 7;
  }

  const now = Math.floor(Date.now() / 1000); // Current time in seconds
  const maxAge = exp - now;

  // Ensure positive maxAge, minimum 60 seconds
  return Math.max(maxAge, 60);
}
