import { LexicalComposer } from "@lexical/react/LexicalComposer";
import { RichTextPlugin } from "@lexical/react/LexicalRichTextPlugin";
import { ContentEditable } from "@lexical/react/LexicalContentEditable";
import { HistoryPlugin } from "@lexical/react/LexicalHistoryPlugin";
import { LexicalErrorBoundary } from "@lexical/react/LexicalErrorBoundary";
import { HeadingNode, QuoteNode } from "@lexical/rich-text";
import { ListNode, ListItemNode } from "@lexical/list";
import { OnChangePlugin } from "@lexical/react/LexicalOnChangePlugin";
import { useLexicalComposerContext } from "@lexical/react/LexicalComposerContext";
import { $getRoot, $insertNodes } from "lexical";
import { $generateNodesFromDOM, $generateHtmlFromNodes } from "@lexical/html";
import { useEffect } from "react";
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

interface LexicalEditorFieldProps {
  value?: string;
  onChange: (html: string) => void;
  placeholder?: string;
}

function InitialContentPlugin({ content }: { content?: string }) {
  const [editor] = useLexicalComposerContext();

  useEffect(() => {
    if (content && content.trim()) {
      editor.update(() => {
        const parser = new DOMParser();
        const dom = parser.parseFromString(content, "text/html");
        const nodes = $generateNodesFromDOM(editor, dom);
        const root = $getRoot();
        root.clear();
        root.select();
        $insertNodes(nodes);
      });
    }
  }, []); // Only run once on mount

  return null;
}

export default function LexicalEditorField({
  value,
  onChange,
  placeholder = "Enter content...",
}: LexicalEditorFieldProps) {
  const initialConfig = {
    namespace: "LexicalEditorField",
    theme,
    onError,
    nodes: [HeadingNode, QuoteNode, ListNode, ListItemNode],
  };

  const handleChange = (editorState: any, editor: any) => {
    editorState.read(() => {
      const htmlString = $generateHtmlFromNodes(editor);
      onChange(htmlString);
    });
  };

  return (
    <LexicalComposer initialConfig={initialConfig}>
      <div className="relative rounded-lg border border-gray-300 bg-white">
        <ToolbarPlugin />
        <div className="relative">
          <RichTextPlugin
            contentEditable={
              <ContentEditable className="min-h-[300px] resize-none p-4 text-base outline-none" />
            }
            placeholder={
              <div className="pointer-events-none absolute left-4 top-4 text-gray-400">
                {placeholder}
              </div>
            }
            ErrorBoundary={LexicalErrorBoundary}
          />
        </div>
        <HistoryPlugin />
        <OnChangePlugin onChange={handleChange} />
        <InitialContentPlugin content={value} />
      </div>
    </LexicalComposer>
  );
}
