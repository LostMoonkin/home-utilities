"use client";
import { useEffect, useState } from "react";
import { ConfigFileInfo, getConfigs, APIResponse } from "../lib/api";
import { cn } from "../lib/utils";
import { FileText, RefreshCw, Plus } from "lucide-react";

interface ConfigListProps {
  onSelect: (name: string) => void;
  selectedName?: string;
  refreshTrigger: number;
  onCreate: () => void;
}

export default function ConfigList({ onSelect, selectedName, refreshTrigger, onCreate }: ConfigListProps) {
  const [configs, setConfigs] = useState<ConfigFileInfo[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchConfigs = async () => {
    setLoading(true);
    try {
      const res = await getConfigs();
      if (res.biz_code === 0 && res.data) {
        setConfigs(res.data);
      }
    } catch (error) {
      console.error("Failed to fetch configs", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfigs();
  }, [refreshTrigger]);

  return (
    <div className="w-64 border-r border-gray-200 h-full flex flex-col bg-gray-50 flex-shrink-0 font-sans">
      <div className="h-14 p-4 border-b border-gray-200 flex justify-between items-center bg-gray-100/50">
        <h2 className="font-semibold text-xs uppercase tracking-wider text-gray-500">Explorer</h2>
        <div className="flex gap-1">
          <button onClick={onCreate} className="p-1.5 hover:bg-gray-200 rounded text-gray-600 transition-colors" title="New Config">
            <Plus size={16} />
          </button>
          <button onClick={fetchConfigs} className="p-1.5 hover:bg-gray-200 rounded text-gray-600 transition-colors" title="Refresh">
            <RefreshCw size={14} className={cn(loading && "animate-spin")} />
          </button>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto pt-2 pb-2">
        {configs.map((file) => (
          <button
            key={file.name}
            onClick={() => onSelect(file.name)}
            className={cn(
              "w-full text-left px-4 py-2 text-sm flex items-center gap-3 transition-colors border-l-2",
              selectedName === file.name
                ? "bg-white border-blue-500 text-blue-600 font-medium"
                : "border-transparent hover:bg-gray-100 text-gray-600 hover:text-gray-900"
            )}
          >
            <FileText size={16} className={cn("flex-shrink-0", selectedName === file.name ? "text-blue-500" : "text-gray-400")} />
            <span className="truncate">{file.name}</span>
          </button>
        ))}
        {configs.length === 0 && !loading && (
          <div className="text-xs text-center text-gray-400 py-8 italic">No configs found</div>
        )}
      </div>
    </div>
  );
}
