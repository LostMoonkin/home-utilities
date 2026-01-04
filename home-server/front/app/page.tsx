"use client";
import { useState } from "react";
import ConfigList from "../components/ConfigList";
import ConfigEditor from "../components/ConfigEditor";
import CreateConfigModal from "../components/CreateConfigModal";
import WelcomeScreen from "../components/WelcomeScreen";
import { applyChanges } from "../lib/api";
import { toast } from "sonner";
import { Play } from "lucide-react";

export default function Home() {
  const [selectedFile, setSelectedFile] = useState<string | undefined>();
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [applying, setApplying] = useState(false);

  const handleApply = async () => {
    setApplying(true);
    try {
        const res = await applyChanges();
        if (res.biz_code === 0) {
            toast.success("Changes applied successfully (Nginx reloaded)");
        } else {
            toast.error(res.message || "Failed to apply changes");
        }
    } catch (error) {
        toast.error("Error applying changes");
    } finally {
        setApplying(false);
    }
  };

  const handleRefresh = () => setRefreshTrigger(p => p + 1);
  const handleCreate = () => setShowCreateModal(true);

  return (
    <div className="flex h-screen w-screen bg-white text-gray-900 overflow-hidden font-sans">
      <ConfigList
        onSelect={setSelectedFile}
        selectedName={selectedFile}
        refreshTrigger={refreshTrigger}
        onCreate={handleCreate}
      />
      
      <main className="flex-1 flex flex-col min-w-0 bg-white relative">
        {/* Top Navigation / App Bar */}
        <div className="h-14 border-b border-gray-200 flex items-center justify-between px-6 bg-white shrink-0">
             <div className="flex items-center gap-3">
                {/* Possibly a logo or just a title if no config selected */}
                {!selectedFile && (
                    <div className="flex items-center gap-2">
                        <div className="h-6 w-6 bg-blue-600 rounded flex items-center justify-center text-white font-bold text-xs">NG</div>
                        <h1 className="font-semibold text-lg tracking-tight text-gray-700">Nginx Manager</h1>
                    </div>
                )}
             </div>
             <div>
                <button
                    onClick={handleApply}
                    disabled={applying}
                    className="flex items-center gap-2 bg-green-50 text-green-700 border border-green-200 px-4 py-1.5 rounded text-sm hover:bg-green-100 hover:border-green-300 transition-all disabled:opacity-70 font-medium"
                    title="Restart Nginx to apply all changes"
                >
                    <Play size={14} className={applying ? "animate-spin" : "fill-current"} />
                    {applying ? "Applying..." : "Apply Changes"}
                </button>
             </div>
        </div>

        <div className="flex-1 overflow-hidden relative flex flex-col">
            {selectedFile ? (
                <ConfigEditor
                    fileName={selectedFile}
                    onSuccess={handleRefresh}
                />
            ) : (
                <WelcomeScreen onCreate={handleCreate} onRefresh={handleRefresh} />
            )}
        </div>
      </main>

      {showCreateModal && (
        <CreateConfigModal
            onClose={() => setShowCreateModal(false)}
            onSuccess={(name) => {
                setShowCreateModal(false);
                setSelectedFile(name);
                handleRefresh();
            }}
        />
      )}
    </div>
  );
}
