import React, { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { 
  Play, Square, RotateCcw, Terminal, HardDrive, Settings, 
  Users, Database, Calendar, Wifi, Download, Package,
  FileText, Folder, Upload, Trash2, Edit, Eye, User,
  Shield, Sword, Diamond, Heart, Clock, MapPin
} from 'lucide-react'

interface ServerDetails {
  id: string
  name: string
  game: string
  gameVersion: string
  status: 'running' | 'stopped' | 'starting' | 'stopping' | 'error'
  node: { name: string; location: string }
  resources: {
    memory: { used: number; total: number }
    cpu: number
    disk: { used: number; total: number }
  }
  players: { online: number; max: number; list: Player[] }
  port: number
  uptime: string
}

interface Player {
  uuid: string
  name: string
  level: number
  health: number
  hunger: number
  gamemode: string
  location: { x: number; y: number; z: number; world: string }
  inventory: InventoryItem[]
  playtime: string
  lastSeen: string
}

interface InventoryItem {
  slot: number
  item: string
  count: number
  enchantments?: string[]
}

interface Plugin {
  name: string
  version: string
  description: string
  author: string
  enabled: boolean
  downloadUrl?: string
}

interface ConsoleLog {
  timestamp: string
  level: 'INFO' | 'WARN' | 'ERROR' | 'DEBUG'
  message: string
}

export default function ServerDetailPage() {
  const { id } = useParams()
  const [server, setServer] = useState<ServerDetails | null>(null)
  const [activeTab, setActiveTab] = useState('console')
  const [consoleLogs, setConsoleLogs] = useState<ConsoleLog[]>([])
  const [consoleInput, setConsoleInput] = useState('')
  const [plugins, setPlugins] = useState<Plugin[]>([])
  const [availablePlugins, setAvailablePlugins] = useState<Plugin[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchServerDetails()
    fetchConsoleLogs()
    fetchPlugins()
    fetchAvailablePlugins()
  }, [id])

  const fetchServerDetails = async () => {
    try {
      // TODO: Replace with actual API call
      setLoading(false)
      // Mock data for now
      setServer(null)
    } catch (error) {
      console.error('Failed to fetch server details:', error)
      setLoading(false)
    }
  }

  const fetchConsoleLogs = async () => {
    try {
      // TODO: Replace with actual API call - WebSocket connection for real-time logs
      setConsoleLogs([])
    } catch (error) {
      console.error('Failed to fetch console logs:', error)
    }
  }

  const fetchPlugins = async () => {
    try {
      // TODO: Replace with actual API call
      setPlugins([])
    } catch (error) {
      console.error('Failed to fetch plugins:', error)
    }
  }

  const fetchAvailablePlugins = async () => {
    try {
      // TODO: Replace with actual API call to plugin repositories
      setAvailablePlugins([])
    } catch (error) {
      console.error('Failed to fetch available plugins:', error)
    }
  }

  const handleServerAction = async (action: 'start' | 'stop' | 'restart') => {
    try {
      // TODO: Replace with actual API call
      console.log(`${action} server ${id}`)
    } catch (error) {
      console.error(`Failed to ${action} server:`, error)
    }
  }

  const handleConsoleCommand = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!consoleInput.trim()) return

    try {
      // TODO: Replace with actual API call
      console.log('Sending command:', consoleInput)
      setConsoleInput('')
    } catch (error) {
      console.error('Failed to send command:', error)
    }
  }

  const handlePluginToggle = async (pluginName: string, enabled: boolean) => {
    try {
      // TODO: Replace with actual API call
      console.log(`${enabled ? 'Enable' : 'Disable'} plugin:`, pluginName)
    } catch (error) {
      console.error('Failed to toggle plugin:', error)
    }
  }

  const handlePluginDownload = async (plugin: Plugin) => {
    try {
      // TODO: Replace with actual API call to download and install plugin
      console.log('Downloading plugin:', plugin.name)
    } catch (error) {
      console.error('Failed to download plugin:', error)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-muted-foreground">Loading server details...</div>
      </div>
    )
  }

  if (!server) {
    return (
      <div className="text-center py-12">
        <div className="text-muted-foreground">Server not found</div>
      </div>
    )
  }

  const tabs = [
    { id: 'console', label: 'Console', icon: Terminal },
    { id: 'files', label: 'Files', icon: Folder },
    { id: 'databases', label: 'Databases', icon: Database },
    { id: 'schedules', label: 'Schedules', icon: Calendar },
    { id: 'users', label: 'Players', icon: Users },
    { id: 'backups', label: 'Backups', icon: HardDrive },
    { id: 'network', label: 'Network', icon: Wifi },
    { id: 'startup', label: 'Startup', icon: Play },
    { id: 'settings', label: 'Settings', icon: Settings },
    { id: 'plugins', label: 'Plugins', icon: Package },
    { id: 'activity', label: 'Activity', icon: Clock }
  ]

  return (
    <div className="space-y-6">
      {/* Server Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">{server.name}</h1>
          <p className="text-muted-foreground">
            {server.game} {server.gameVersion} • {server.node.name} ({server.node.location})
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span className={`px-3 py-1 rounded-full text-xs ${
            server.status === 'running' ? 'bg-green-100 text-green-800' :
            server.status === 'stopped' ? 'bg-gray-100 text-gray-800' :
            server.status === 'starting' ? 'bg-blue-100 text-blue-800' :
            server.status === 'stopping' ? 'bg-yellow-100 text-yellow-800' :
            'bg-red-100 text-red-800'
          }`}>
            {server.status}
          </span>
          <button 
            onClick={() => handleServerAction('start')}
            disabled={server.status === 'running' || server.status === 'starting'}
            className="p-2 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:opacity-50"
          >
            <Play className="w-5 h-5" />
          </button>
          <button 
            onClick={() => handleServerAction('stop')}
            disabled={server.status === 'stopped' || server.status === 'stopping'}
            className="p-2 bg-red-500 text-white rounded-lg hover:bg-red-600 disabled:opacity-50"
          >
            <Square className="w-5 h-5" />
          </button>
          <button 
            onClick={() => handleServerAction('restart')}
            disabled={server.status !== 'running'}
            className="p-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 disabled:opacity-50"
          >
            <RotateCcw className="w-5 h-5" />
          </button>
        </div>
      </div>

      {/* Resource Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-card border rounded-lg p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-muted-foreground">Memory</span>
            <span className="text-sm font-medium">
              {Math.round(server.resources.memory.used / 1024)}GB / {Math.round(server.resources.memory.total / 1024)}GB
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-blue-500 h-2 rounded-full" 
              style={{ width: `${(server.resources.memory.used / server.resources.memory.total) * 100}%` }}
            ></div>
          </div>
        </div>

        <div className="bg-card border rounded-lg p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-muted-foreground">CPU</span>
            <span className="text-sm font-medium">{server.resources.cpu}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-green-500 h-2 rounded-full" 
              style={{ width: `${server.resources.cpu}%` }}
            ></div>
          </div>
        </div>

        <div className="bg-card border rounded-lg p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-muted-foreground">Disk</span>
            <span className="text-sm font-medium">
              {Math.round(server.resources.disk.used / 1024)}GB / {Math.round(server.resources.disk.total / 1024)}GB
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-purple-500 h-2 rounded-full" 
              style={{ width: `${(server.resources.disk.used / server.resources.disk.total) * 100}%` }}
            ></div>
          </div>
        </div>

        <div className="bg-card border rounded-lg p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-muted-foreground">Players</span>
            <span className="text-sm font-medium">{server.players.online} / {server.players.max}</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-orange-500 h-2 rounded-full" 
              style={{ width: `${(server.players.online / server.players.max) * 100}%` }}
            ></div>
          </div>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="border-b">
        <nav className="flex space-x-8 overflow-x-auto">
          {tabs.map((tab) => {
            const Icon = tab.icon
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-muted-foreground hover:text-foreground hover:border-gray-300'
                }`}
              >
                <Icon className="w-4 h-4" />
                {tab.label}
              </button>
            )
          })}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="min-h-[500px]">
        {activeTab === 'console' && <ConsoleTab 
          logs={consoleLogs} 
          onCommand={handleConsoleCommand}
          input={consoleInput}
          setInput={setConsoleInput}
        />}
        {activeTab === 'files' && <FilesTab serverId={id!} />}
        {activeTab === 'users' && <PlayersTab players={server.players.list} />}
        {activeTab === 'plugins' && <PluginsTab 
          plugins={plugins}
          availablePlugins={availablePlugins}
          onToggle={handlePluginToggle}
          onDownload={handlePluginDownload}
        />}
        {activeTab === 'settings' && <SettingsTab server={server} />}
        {/* Add other tab components as needed */}
      </div>
    </div>
  )
}

// Console Tab Component
function ConsoleTab({ logs, onCommand, input, setInput }: {
  logs: ConsoleLog[]
  onCommand: (e: React.FormEvent) => void
  input: string
  setInput: (value: string) => void
}) {
  return (
    <div className="bg-card border rounded-lg p-6">
      <div className="bg-black rounded-lg p-4 h-96 font-mono text-sm overflow-auto">
        {logs.length === 0 ? (
          <div className="text-gray-500">No console logs available. Start the server to see logs.</div>
        ) : (
          logs.map((log, index) => (
            <div key={index} className={`mb-1 ${
              log.level === 'ERROR' ? 'text-red-400' :
              log.level === 'WARN' ? 'text-yellow-400' :
              log.level === 'DEBUG' ? 'text-gray-400' :
              'text-green-400'
            }`}>
              <span className="text-gray-500">[{log.timestamp}]</span> 
              <span className="text-blue-400">[{log.level}]</span> {log.message}
            </div>
          ))
        )}
      </div>
      <form onSubmit={onCommand} className="mt-4">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Enter command..."
          className="w-full px-4 py-2 border rounded-lg font-mono"
        />
      </form>
    </div>
  )
}

// Players Tab Component
function PlayersTab({ players }: { players: Player[] }) {
  const [selectedPlayer, setSelectedPlayer] = useState<Player | null>(null)

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-card border rounded-lg p-6">
          <h3 className="font-semibold mb-4">Online Players ({players.length})</h3>
          <div className="space-y-2">
            {players.length === 0 ? (
              <div className="text-muted-foreground text-center py-8">No players online</div>
            ) : (
              players.map((player) => (
                <div 
                  key={player.uuid}
                  onClick={() => setSelectedPlayer(player)}
                  className="flex items-center gap-3 p-3 rounded-lg hover:bg-muted cursor-pointer"
                >
                  <div className="w-8 h-8 bg-blue-500 rounded flex items-center justify-center">
                    <User className="w-4 h-4 text-white" />
                  </div>
                  <div className="flex-1">
                    <div className="font-medium">{player.name}</div>
                    <div className="text-sm text-muted-foreground">Level {player.level}</div>
                  </div>
                  <div className="text-right text-sm">
                    <div className="flex items-center gap-1">
                      <Heart className="w-3 h-3 text-red-500" />
                      {player.health}/20
                    </div>
                    <div className="text-muted-foreground">{player.gamemode}</div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {selectedPlayer && (
          <div className="bg-card border rounded-lg p-6">
            <h3 className="font-semibold mb-4">Player Details</h3>
            <div className="space-y-4">
              <div>
                <h4 className="font-medium mb-2">Basic Info</h4>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground">Name:</span> {selectedPlayer.name}
                  </div>
                  <div>
                    <span className="text-muted-foreground">Level:</span> {selectedPlayer.level}
                  </div>
                  <div>
                    <span className="text-muted-foreground">Health:</span> {selectedPlayer.health}/20
                  </div>
                  <div>
                    <span className="text-muted-foreground">Hunger:</span> {selectedPlayer.hunger}/20
                  </div>
                  <div>
                    <span className="text-muted-foreground">Gamemode:</span> {selectedPlayer.gamemode}
                  </div>
                  <div>
                    <span className="text-muted-foreground">Playtime:</span> {selectedPlayer.playtime}
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-medium mb-2">Location</h4>
                <div className="text-sm">
                  <div className="flex items-center gap-1">
                    <MapPin className="w-3 h-3" />
                    {selectedPlayer.location.world}: {selectedPlayer.location.x}, {selectedPlayer.location.y}, {selectedPlayer.location.z}
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-medium mb-2">Inventory</h4>
                <div className="grid grid-cols-9 gap-1">
                  {Array.from({ length: 36 }, (_, i) => {
                    const item = selectedPlayer.inventory.find(item => item.slot === i)
                    return (
                      <div key={i} className="aspect-square border rounded bg-muted/50 p-1 text-xs flex items-center justify-center">
                        {item ? (
                          <div className="text-center">
                            <div className="font-mono">{item.item.split(':')[1]?.slice(0, 3) || '?'}</div>
                            <div className="text-[10px]">{item.count}</div>
                          </div>
                        ) : (
                          <div className="w-full h-full bg-gray-100 rounded"></div>
                        )}
                      </div>
                    )
                  })}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

// Plugins Tab Component
function PluginsTab({ plugins, availablePlugins, onToggle, onDownload }: {
  plugins: Plugin[]
  availablePlugins: Plugin[]
  onToggle: (name: string, enabled: boolean) => void
  onDownload: (plugin: Plugin) => void
}) {
  const [activePluginTab, setActivePluginTab] = useState<'installed' | 'available'>('installed')

  return (
    <div className="space-y-6">
      <div className="flex gap-4">
        <button
          onClick={() => setActivePluginTab('installed')}
          className={`px-4 py-2 rounded-lg ${
            activePluginTab === 'installed' ? 'bg-blue-500 text-white' : 'bg-muted'
          }`}
        >
          Installed Plugins ({plugins.length})
        </button>
        <button
          onClick={() => setActivePluginTab('available')}
          className={`px-4 py-2 rounded-lg ${
            activePluginTab === 'available' ? 'bg-blue-500 text-white' : 'bg-muted'
          }`}
        >
          Available Plugins ({availablePlugins.length})
        </button>
      </div>

      {activePluginTab === 'installed' ? (
        <div className="grid gap-4">
          {plugins.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              No plugins installed
            </div>
          ) : (
            plugins.map((plugin) => (
              <div key={plugin.name} className="bg-card border rounded-lg p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-semibold">{plugin.name}</h4>
                    <p className="text-sm text-muted-foreground">{plugin.description}</p>
                    <p className="text-xs text-muted-foreground">v{plugin.version} by {plugin.author}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => onToggle(plugin.name, !plugin.enabled)}
                      className={`px-3 py-1 rounded text-sm ${
                        plugin.enabled 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-gray-100 text-gray-800'
                      }`}
                    >
                      {plugin.enabled ? 'Enabled' : 'Disabled'}
                    </button>
                    <button className="p-2 hover:bg-muted rounded">
                      <Settings className="w-4 h-4" />
                    </button>
                    <button className="p-2 hover:bg-muted rounded text-red-500">
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      ) : (
        <div className="grid gap-4">
          {availablePlugins.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              Loading available plugins...
            </div>
          ) : (
            availablePlugins.map((plugin) => (
              <div key={plugin.name} className="bg-card border rounded-lg p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-semibold">{plugin.name}</h4>
                    <p className="text-sm text-muted-foreground">{plugin.description}</p>
                    <p className="text-xs text-muted-foreground">v{plugin.version} by {plugin.author}</p>
                  </div>
                  <button
                    onClick={() => onDownload(plugin)}
                    className="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                  >
                    <Download className="w-4 h-4" />
                    Download
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  )
}

// Files Tab Component
function FilesTab({ serverId }: { serverId: string }) {
  return (
    <div className="bg-card border rounded-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-semibold">File Manager</h3>
        <div className="flex gap-2">
          <button className="flex items-center gap-2 px-3 py-1 bg-blue-500 text-white rounded">
            <Upload className="w-4 h-4" />
            Upload
          </button>
          <button className="flex items-center gap-2 px-3 py-1 bg-green-500 text-white rounded">
            <FileText className="w-4 h-4" />
            New File
          </button>
        </div>
      </div>
      <div className="text-center py-12 text-muted-foreground">
        File manager will be implemented with real server file system access
      </div>
    </div>
  )
}

// Settings Tab Component
function SettingsTab({ server }: { server: ServerDetails }) {
  return (
    <div className="bg-card border rounded-lg p-6">
      <h3 className="font-semibold mb-4">Server Settings</h3>
      <div className="text-center py-12 text-muted-foreground">
        Server settings panel will be implemented
      </div>
    </div>
  )
}
