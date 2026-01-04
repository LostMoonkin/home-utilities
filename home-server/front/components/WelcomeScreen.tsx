"use client";
import { FileText, Plus, RefreshCw, Edit3 } from "lucide-react";

interface WelcomeScreenProps {
  onCreate: () => void;
  onRefresh: () => void;
}

export default function WelcomeScreen({ onCreate, onRefresh }: WelcomeScreenProps) {
  return (
    <div className="flex-1 overflow-y-auto bg-gray-50 p-8 font-sans text-gray-800">
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Header Section */}
        <section className="space-y-4">
          <div className="flex items-center gap-3 text-gray-500 mb-2">
            <Edit3 size={24} />
            <span className="text-lg">Nginx Config Editor</span>
          </div>
          <h1 className="text-3xl font-light text-gray-700">Welcome to Nginx Manager</h1>
          <p className="text-gray-600 max-w-2xl leading-relaxed">
            Nginx Manager is a lightweight tool that allows you to create, open, and edit Nginx configuration files directly on your server.
            <br />
            To get started, select an option below or choose a file from the sidebar.
          </p>
        </section>

        {/* Action Buttons */}
        <div className="flex flex-wrap gap-4">
          <button 
            onClick={onRefresh}
            className="flex items-center gap-2 bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded shadow-sm transition-colors font-medium"
          >
            <RefreshCw size={20} />
            Refresh Config List
          </button>
          
          <button 
            onClick={onCreate}
            className="flex items-center gap-2 bg-sky-500 hover:bg-sky-600 text-white px-6 py-3 rounded shadow-sm transition-colors font-medium"
          >
            <Plus size={20} />
            Create New Config
          </button>
        </div>

        {/* Instructions Section */}
        <section className="space-y-4 pt-4">
          <h2 className="text-xl font-light text-gray-700">How to use Nginx Manager</h2>
          <ol className="list-decimal list-inside space-y-2 text-gray-600 ml-2">
            <li>Select a configuration file from the list on the left sidebar.</li>
            <li>The file contents will be displayed in the editor where you can make changes.</li>
            <li>After editing, press the <span className="font-semibold text-gray-700">"Save Changes"</span> button to update the file on the server.</li>
            <li>Use the <span className="font-semibold text-gray-700">"Apply Changes"</span> button in the top bar to reload Nginx and apply your updates.</li>
          </ol>
        </section>

        {/* Features Section */}
        <section className="space-y-4 pt-4">
          <h2 className="text-xl font-light text-gray-700">Features</h2>
          <ul className="list-disc list-inside space-y-2 text-gray-600 ml-2">
            <li><strong>Syntax Highlighting:</strong> Advanced editor with Nginx syntax support.</li>
            <li><strong>Safe Editing:</strong> Validates file extensions and provides feedback.</li>
            <li><strong>Direct Integration:</strong> Interacts directly with your server's Nginx configuration directory.</li>
            <li><strong>Instant Apply:</strong> Restart Nginx services with a single click.</li>
          </ul>
        </section>
      </div>
    </div>
  );
}
