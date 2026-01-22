import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Search, Server, Play, Square, RotateCcw, Trash2, Settings } from 'lucide-react'
import { cn } from '../../lib/utils'

interface GameServer {
  id: string
  name: string
  game: string
  gameVersion: string
  nodeId: string
  nodeName: string
  status: 'running' | 'stopped' | 'starting' | 'stopping' | 'error'
  resources: {
    memory: { used: number; total: number }
    cpu: number
    disk: { used: number; total: number }
  }
  players: { online: number; max: number }
  port: number
  createdAt: string
  lastStarted?: string
}

export default function ServersPage() {
  const [servers, setServers] = useState<GameServer[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [showCreateModal, setShowCreateModal] = useState(false)

  useEffect(() => {
    fetchServers()
  }, [])

  const fetchServers = async () => {
    try {
      // TODO: Replace with actual API call
      setLoading(false)
      setServers([]) // Start with empty servers
    } catch (error) {
      console.error('Failed to fetch servers:', error)
      setLoading(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running': return 'bg-green-100 text-green-800'
      case 'stopped': return 'bg-gray-100 text-gray-800'
      case 'starting': return 'bg-blue-100 text-blue-800'
      case 'stopping': return 'bg-yellow-100 text-yellow-800'
      case 'error': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const filteredServers = servers.filter(server =>
    server.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    server.game.toLowerCase().includes(searchTerm.toLowerCase())
  )

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-foreground">Servers</h1>
            <p className="text-muted-foreground">Manage your game servers</p>
          </div>
        </div>
        <div className="flex items-center justify-center py-12">
          <div className="text-muted-foreground">Loading servers...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Servers</h1>
          <p className="text-muted-foreground">Manage your game servers</p>
        </div>
        <button 
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
        >
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
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border rounded-lg"
          />
        </div>
      </div>

      {filteredServers.length === 0 ? (
        <div className="text-center py-12">
          <Server className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">No servers found</h3>
          <p className="text-muted-foreground mb-4">
            {searchTerm ? 'No servers match your search criteria' : 'Create your first game server to get started'}
          </p>
          {!searchTerm && (
            <button 
              onClick={() => setShowCreateModal(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 mx-auto"
            >
              <Plus className="h-4 w-4" />
              Create Server
            </button>
          )}
        </div>
      ) : (
        <div className="grid gap-4">
          {filteredServers.map((server) => (
            <div key={server.id} className="bg-card border rounded-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                  <div className="p-3 bg-accent rounded-lg">
                    <Server className="w-6 h-6 text-foreground" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-lg">{server.name}</h3>
                    <p className="text-sm text-muted-foreground">{server.game} {server.gameVersion}</p>
                    <p className="text-xs text-muted-foreground">Node: {server.nodeName} | Port: {server.port}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`px-2 py-1 rounded-full text-xs ${getStatusColor(server.status)}`}>
                    {server.status}
                  </span>
                  <div className="flex gap-1">
                    <button className="p-2 hover:bg-muted rounded-lg text-green-600">
                      <Play className="h-4 w-4" />
                    </button>
                    <button className="p-2 hover:bg-muted rounded-lg text-red-600">
                      <Square className="h-4 w-4" />
                    </button>
                    <button className="p-2 hover:bg-muted rounded-lg text-blue-600">
                      <RotateCcw className="h-4 w-4" />
                    </button>
                    <Link to={`/servers/${server.id}`} className="p-2 hover:bg-muted rounded-lg">
                      <Settings className="h-4 w-4" />
                    </Link>
                    <button className="p-2 hover:bg-muted rounded-lg text-red-500">
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-4 gap-4">
                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="text-sm font-medium mb-1">Memory</div>
                  <div className="text-lg font-semibold">
                    {Math.round(server.resources.memory.used / 1024)}GB / {Math.round(server.resources.memory.total / 1024)}GB
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div 
                      className="bg-blue-500 h-2 rounded-full" 
                      style={{ width: `${(server.resources.memory.used / server.resources.memory.total) * 100}%` }}
                    ></div>
                  </div>
                </div>

                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="text-sm font-medium mb-1">CPU</div>
                  <div className="text-lg font-semibold">{server.resources.cpu}%</div>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div 
                      className="bg-green-500 h-2 rounded-full" 
                      style={{ width: `${server.resources.cpu}%` }}
                    ></div>
                  </div>
                </div>

                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="text-sm font-medium mb-1">Disk</div>
                  <div className="text-lg font-semibold">
                    {Math.round(server.resources.disk.used / 1024)}GB / {Math.round(server.resources.disk.total / 1024)}GB
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div 
                      className="bg-purple-500 h-2 rounded-full" 
                      style={{ width: `${(server.resources.disk.used / server.resources.disk.total) * 100}%` }}
                    ></div>
                  </div>
                </div>

                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="text-sm font-medium mb-1">Players</div>
                  <div className="text-lg font-semibold">
                    {server.players.online} / {server.players.max}
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div 
                      className="bg-orange-500 h-2 rounded-full" 
                      style={{ width: `${(server.players.online / server.players.max) * 100}%` }}
                    ></div>
                  </div>
                </div>
              </div>

              <div className="mt-4 pt-4 border-t text-xs text-muted-foreground">
                Created: {server.createdAt} {server.lastStarted && `| Last started: ${server.lastStarted}`}
              </div>
            </div>
          ))}
        </div>
      )}

      {showCreateModal && <CreateServerModal onClose={() => setShowCreateModal(false)} onCreated={fetchServers} />}
    </div>
  )
}

function CreateServerModal({ onClose, onCreated }: { onClose: () => void; onCreated: () => void }) {
  const [formData, setFormData] = useState({
    name: '',
    game: 'minecraft-java',
    gameVersion: 'latest',
    nodeId: '',
    memory: '2048',
    maxPlayers: '20',
    port: '25565'
  })
  const [loading, setLoading] = useState(false)
  const [nodes] = useState([]) // TODO: Fetch available nodes

  const gameOptions = [
    { value: 'minecraft-java', label: 'Minecraft Java Edition' },
    { value: 'minecraft-bedrock', label: 'Minecraft Bedrock Edition' },
    { value: 'rust', label: 'Rust' },
    { value: 'csgo', label: 'Counter-Strike: Global Offensive' },
    { value: 'gmod', label: 'Garry\'s Mod' },
    { value: 'ark', label: 'ARK: Survival Evolved' }
  ]

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    
    try {
      // TODO: Replace with actual API call
      console.log('Creating server:', formData)
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      onCreated()
      onClose()
    } catch (error) {
      console.error('Failed to create server:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-card border rounded-lg p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
        <h2 className="text-xl font-semibold mb-4">Create New Server</h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Server Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              placeholder="e.g., My Minecraft Server"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Game</label>
            <select
              value={formData.game}
              onChange={(e) => setFormData({ ...formData, game: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              required
            >
              {gameOptions.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Game Version</label>
            <input
              type="text"
              value={formData.gameVersion}
              onChange={(e) => setFormData({ ...formData, gameVersion: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              placeholder="latest"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Node</label>
            <select
              value={formData.nodeId}
              onChange={(e) => setFormData({ ...formData, nodeId: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              required
            >
              <option value="">Select a node</option>
              {nodes.map((node: any) => (
                <option key={node.id} value={node.id}>
                  {node.name} ({node.location})
                </option>
              ))}
            </select>
            {nodes.length === 0 && (
              <p className="text-xs text-red-500 mt-1">No nodes available. Please add a node first.</p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Memory (MB)</label>
              <input
                type="number"
                value={formData.memory}
                onChange={(e) => setFormData({ ...formData, memory: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg"
                min="512"
                max="16384"
                step="256"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Max Players</label>
              <input
                type="number"
                value={formData.maxPlayers}
                onChange={(e) => setFormData({ ...formData, maxPlayers: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg"
                min="1"
                max="100"
                required
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Port</label>
            <input
              type="number"
              value={formData.port}
              onChange={(e) => setFormData({ ...formData, port: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              min="1024"
              max="65535"
              required
            />
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border rounded-lg hover:bg-muted"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading || nodes.length === 0}
              className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg disabled:opacity-50"
            >
              {loading ? 'Creating...' : 'Create Server'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
