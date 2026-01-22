import React from 'react'
import { Link } from 'react-router-dom'
import { Plus, Search, Server } from 'lucide-react'
import { cn } from '../../lib/utils'

const servers = [
  { id: '1', name: 'Minecraft SMP', game: 'Minecraft Java', status: 'running', memory: 75, cpu: 45 },
  { id: '2', name: 'Creative Server', game: 'Minecraft Java', status: 'stopped', memory: 0, cpu: 0 },
  { id: '3', name: 'Rust Survival', game: 'Rust', status: 'running', memory: 82, cpu: 60 },
]

export default function ServersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Servers</h1>
          <p className="text-muted-foreground">Manage your game servers</p>
        </div>
        <button className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-lg hover:bg-primary/90">
          <Plus className="w-4 h-4" />
          New Server
        </button>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search servers..."
            className="w-full pl-10 pr-4 py-2 bg-input border border-border rounded-lg text-foreground"
          />
        </div>
      </div>

      <div className="grid gap-4">
        {servers.map((server) => (
          <Link
            key={server.id}
            to={`/servers/${server.id}`}
            className="block bg-card border border-border rounded-xl p-6 hover:border-primary/50 transition-colors"
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                <div className="p-3 bg-accent rounded-lg">
                  <Server className="w-6 h-6 text-foreground" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">{server.name}</h3>
                  <p className="text-sm text-muted-foreground">{server.game}</p>
                </div>
              </div>
              <div className="flex items-center gap-6">
                <div className="text-right">
                  <p className="text-sm text-muted-foreground">Memory</p>
                  <p className="font-medium text-foreground">{server.memory}%</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-muted-foreground">CPU</p>
                  <p className="font-medium text-foreground">{server.cpu}%</p>
                </div>
                <div className={cn(
                  'px-3 py-1 rounded-full text-xs font-medium',
                  server.status === 'running' ? 'bg-green-500/20 text-green-500' : 'bg-gray-500/20 text-gray-500'
                )}>
                  {server.status}
                </div>
              </div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
