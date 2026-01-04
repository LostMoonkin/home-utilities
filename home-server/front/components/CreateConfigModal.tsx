"use client";
import { useState } from "react";
import { createConfig } from "../lib/api";
import { toast } from "sonner";
import { X, Loader2 } from "lucide-react";

interface CreateConfigModalProps {
  onClose: () => void;
  onSuccess: (name: string) => void;
}

export default function CreateConfigModal({ onClose, onSuccess }: CreateConfigModalProps) {
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.endsWith(".conf")) {
      toast.error("Config name must end with .conf");
      return;
    }
    setLoading(true);
    try {
      // Create empty config or default basic content
      const defaultContent = btoa("# New config\n"); 
      const res = await createConfig(name, defaultContent);
      if (res.biz_code === 0) {
        toast.success("Config created successfully");
        onSuccess(name);
      } else {
        toast.error(res.message || "Failed to create config");
      }
    } catch (error) {
      toast.error("Error creating config");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-[400px] p-6 relative">
        <button onClick={onClose} className="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
            <X size={20} />
        </button>
        <h2 className="text-xl font-semibold mb-4">New Configuration</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">File Name</label>
                <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="example.conf"
                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    autoFocus
                />
                <p className="text-xs text-gray-500 mt-1">Must end with .conf</p>
            </div>
            <div className="flex justify-end gap-2 pt-2">
                <button type="button" onClick={onClose} className="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-md">Cancel</button>
                <button
                    type="submit"
                    disabled={loading || !name}
                    className="flex items-center gap-2 bg-black text-white px-4 py-2 rounded-md text-sm hover:bg-gray-800 disabled:opacity-50"
                >
                    {loading && <Loader2 size={14} className="animate-spin" />}
                    Create
                </button>
            </div>
        </form>
      </div>
    </div>
  );
}
