import { LexicalComposer } from "@lexical/react/LexicalComposer";
import { RichTextPlugin } from "@lexical/react/LexicalRichTextPlugin";
import { ContentEditable } from "@lexical/react/LexicalContentEditable";
import { HistoryPlugin } from "@lexical/react/LexicalHistoryPlugin";
import { LexicalErrorBoundary } from "@lexical/react/LexicalErrorBoundary";
import { AutoFocusPlugin } from "@lexical/react/LexicalAutoFocusPlugin";
import { HeadingNode, QuoteNode } from "@lexical/rich-text";
import { ListNode, ListItemNode } from "@lexical/list";
import ToolbarPlugin from "./ToolbarPlugin";

const theme = {
  paragraph: "mb-2 text-base",
  text: {
    bold: "font-bold",
    italic: "italic",
    underline: "underline",
    strikethrough: "line-through",
    code: "bg-gray-100 px-1 py-0.5 rounded font-mono text-sm",
  },
  heading: {
    h1: "text-3xl font-bold mb-4",
    h2: "text-2xl font-bold mb-3",
    h3: "text-xl font-bold mb-2",
  },
  quote: "border-l-4 border-gray-300 pl-4 italic text-gray-700",
  list: {
    ul: "list-disc list-inside ml-4 mb-2",
    ol: "list-decimal list-inside ml-4 mb-2",
    listitem: "mb-1",
  },
};

function onError(error: Error) {
  console.error(error);
}

export default function LexicalEditor() {
  const initialConfig = {
    namespace: "LexicalEditor",
    theme,
    onError,
    nodes: [HeadingNode, QuoteNode, ListNode, ListItemNode],
  };

  return (
    <LexicalComposer initialConfig={initialConfig}>
      <div className="relative rounded-lg border border-gray-300 bg-white">
        <ToolbarPlugin />
        <div className="relative">
          <RichTextPlugin
            contentEditable={
              <ContentEditable className="min-h-[450px] resize-none p-4 text-base outline-none" />
            }
            ErrorBoundary={LexicalErrorBoundary}
          />
        </div>
        <HistoryPlugin />
        <AutoFocusPlugin />
      </div>
    </LexicalComposer>
  );
}
