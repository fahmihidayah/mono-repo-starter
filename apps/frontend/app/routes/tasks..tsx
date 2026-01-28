import type { Route } from "../+types/root";

// Stub tasks route - replace wit\
export async function loader({ params }: Route.LoaderArgs) {
  // TODO: Implement with Go API
  return { tasks: [] };
}

export default function TasksPage() {
  return (
    <div>
      <h1>Tasks</h1>
      <p>TODO: Implement tasks page</p>
    </div>
  );
}
