"use client";
import { useEffect, useState, useRef } from "react";
import Editor from "@monaco-editor/react";
import { getConfig, updateConfig, APIResponse } from "../lib/api";
import { toast } from "sonner";
import { Save, Loader2, RotateCcw, FileText } from "lucide-react";

interface ConfigEditorProps {
  fileName?: string;
  onSuccess: () => void;
}

export default function ConfigEditor({ fileName, onSuccess }: ConfigEditorProps) {
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);
  const [originalContent, setOriginalContent] = useState("");
  const [saving, setSaving] = useState(false);
  const editorRef = useRef<any>(null);

  useEffect(() => {
    if (fileName) {
      loadFile(fileName);
    } else {
      setContent("");
      setOriginalContent("");
    }
  }, [fileName]);

  const loadFile = async (name: string) => {
    setLoading(true);
    try {
      const res = await getConfig(name);
      if (res.biz_code === 0 && res.data && res.data[name]) {
        try {
          const decoded = atob(res.data[name]);
          setContent(decoded);
          setOriginalContent(decoded);
        } catch (e) {
          toast.error("Failed to decode content");
        }
      } else {
        toast.error("Failed to load config");
      }
    } catch (error) {
      toast.error("Error loading config");
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!fileName) return;
    setSaving(true);
    try {
      const currentEncoded = btoa(originalContent);
      const expectedEncoded = btoa(content);
      const res = await updateConfig(fileName, currentEncoded, expectedEncoded);
      if (res.biz_code === 0) {
        toast.success("Saved successfully");
        setOriginalContent(content);
        onSuccess();
      } else {
        toast.error(res.message || "Failed to save");
      }
    } catch (error) {
      toast.error("Error saving config");
    } finally {
      setSaving(false);
    }
  };

  if (!fileName) return <div className="flex-1 flex items-center justify-center text-gray-400 text-sm">Select a file to edit or create new</div>;

  if (loading) return <div className="flex-1 flex items-center justify-center bg-gray-50"><Loader2 className="animate-spin text-blue-500" /></div>;

  return (
    <div className="flex flex-col h-full bg-white flex-1 min-w-0">
      <div className="h-14 border-b border-gray-200 flex items-center justify-between px-4 bg-white">
        <div className="flex items-center gap-3">
          <div className="bg-blue-100 p-1.5 rounded text-blue-600">
            <FileText size={16} />
          </div>
          <div className="flex flex-col">
            <span className="font-medium text-gray-800 text-sm leading-tight">{fileName}</span>
            <span className="text-[10px] text-gray-400 uppercase tracking-wider">Nginx Config</span>
          </div>
          {content !== originalContent && (
            <span className="ml-2 text-[10px] text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full border border-amber-200 font-medium tracking-wide uppercase">Unsaved</span>
          )}
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => loadFile(fileName)}
            className="flex items-center gap-2 px-3 py-1.5 rounded text-xs font-medium text-gray-600 hover:bg-gray-100 transition-colors border border-gray-200"
            title="Reset to server version"
          >
            <RotateCcw size={14} />
            <span>Reset</span>
          </button>
          <button
            onClick={handleSave}
            disabled={saving || content === originalContent}
            className="flex items-center gap-2 bg-blue-600 text-white px-4 py-1.5 rounded text-xs font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors shadow-sm"
          >
            {saving ? <Loader2 size={14} className="animate-spin" /> : <Save size={14} />}
            Save Changes
          </button>
        </div>
      </div>
      <div className="flex-1 pt-0 bg-white">
        <Editor
          height="100%"
          defaultLanguage="shell"
          value={content}
          theme="vs-light"
          onChange={(val) => setContent(val || "")}
          onMount={(editor) => editorRef.current = editor}
          options={{
            minimap: { enabled: true },
            scrollBeyondLastLine: false,
            fontSize: 14,
            fontFamily: "'JetBrains Mono', 'Fira Code', Consolas, monospace",
            lineNumbers: "on",
            renderWhitespace: "selection",
            padding: { top: 20, bottom: 20 },
            lineHeight: 24,
          }}
        />
      </div>
    </div>
  );
}
