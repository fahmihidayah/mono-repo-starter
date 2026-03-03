import { PageHeader } from "~/components/layout/table/table-list";
import LexicalEditor from "~/components/editor/LexicalEditor";

export function meta() {
  return [
    { title: "Lexical Editor - Dashboard" },
    { name: "description", content: "Test Lexical rich text editor" },
  ];
}

export default function LexicalDashboardPage() {
  return (
    <div className="flex-1 p-6">
      <div className="space-y-6">
        <PageHeader
          title="Lexical Editor"
          description="A powerful rich text editor built with Lexical"
        />

        <LexicalEditor />

        <div className="rounded-lg border bg-gray-50 p-4">
          <h3 className="mb-2 font-medium">Features:</h3>
          <ul className="list-inside list-disc space-y-1 text-sm text-gray-700">
            <li>Rich text editing with Lexical framework</li>
            <li>Auto-focus on mount</li>
            <li>Undo/Redo support with History plugin</li>
            <li>Error boundary for graceful error handling</li>
            <li>Clean and minimal UI</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
