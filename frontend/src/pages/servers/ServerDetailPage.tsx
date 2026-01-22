import React from 'react'
import { useParams } from 'react-router-dom'
import { Play, Square, RotateCcw, Terminal, HardDrive, Settings } from 'lucide-react'

export default function ServerDetailPage() {
  const { id } = useParams()

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Minecraft SMP</h1>
          <p className="text-muted-foreground">Server ID: {id}</p>
        </div>
        <div className="flex items-center gap-2">
          <button className="p-2 bg-green-500 text-white rounded-lg hover:bg-green-600">
            <Play className="w-5 h-5" />
          </button>
          <button className="p-2 bg-red-500 text-white rounded-lg hover:bg-red-600">
            <Square className="w-5 h-5" />
          </button>
          <button className="p-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600">
            <RotateCcw className="w-5 h-5" />
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-card border border-border rounded-xl p-6">
            <div className="flex items-center gap-2 mb-4">
              <Terminal className="w-5 h-5 text-foreground" />
              <h2 className="font-semibold text-foreground">Console</h2>
            </div>
            <div className="bg-black rounded-lg p-4 h-64 font-mono text-sm text-green-400 overflow-auto">
              <p>[Server] Starting Minecraft server...</p>
              <p>[Server] Loading world...</p>
              <p>[Server] Done! Server is ready.</p>
            </div>
            <input
              type="text"
              placeholder="Enter command..."
              className="w-full mt-4 px-4 py-2 bg-input border border-border rounded-lg text-foreground"
            />
          </div>
        </div>

        <div className="space-y-6">
          <div className="bg-card border border-border rounded-xl p-6">
            <h2 className="font-semibold text-foreground mb-4">Resources</h2>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-muted-foreground">Memory</span>
                  <span className="text-foreground">1.5GB / 2GB</span>
                </div>
                <div className="h-2 bg-accent rounded-full overflow-hidden">
                  <div className="h-full bg-blue-500 w-3/4" />
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-muted-foreground">CPU</span>
                  <span className="text-foreground">45%</span>
                </div>
                <div className="h-2 bg-accent rounded-full overflow-hidden">
                  <div className="h-full bg-green-500 w-[45%]" />
                </div>
              </div>
              <div>
                <div className="flex justify-between text-sm mb-1">
                  <span className="text-muted-foreground">Disk</span>
                  <span className="text-foreground">8GB / 20GB</span>
                </div>
                <div className="h-2 bg-accent rounded-full overflow-hidden">
                  <div className="h-full bg-purple-500 w-[40%]" />
                </div>
              </div>
            </div>
          </div>

          <div className="bg-card border border-border rounded-xl p-6">
            <h2 className="font-semibold text-foreground mb-4">Quick Actions</h2>
            <div className="space-y-2">
              <button className="w-full flex items-center gap-2 px-4 py-2 bg-accent rounded-lg text-foreground hover:bg-accent/80">
                <HardDrive className="w-4 h-4" /> File Manager
              </button>
              <button className="w-full flex items-center gap-2 px-4 py-2 bg-accent rounded-lg text-foreground hover:bg-accent/80">
                <Settings className="w-4 h-4" /> Settings
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
