import React from 'react'
import { HardDrive, Plus } from 'lucide-react'

export default function NodesPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Nodes</h1>
          <p className="text-muted-foreground">Manage your server nodes</p>
        </div>
        <button className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-lg">
          <Plus className="w-4 h-4" /> Add Node
        </button>
      </div>
      <div className="grid gap-4">
        {[1, 2, 3].map((i) => (
          <div key={i} className="bg-card border border-border rounded-xl p-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                <div className="p-3 bg-green-500/20 rounded-lg">
                  <HardDrive className="w-6 h-6 text-green-500" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">Node {i}</h3>
                  <p className="text-sm text-muted-foreground">node{i}.example.com</p>
                </div>
              </div>
              <span className="px-3 py-1 bg-green-500/20 text-green-500 rounded-full text-xs">Online</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
