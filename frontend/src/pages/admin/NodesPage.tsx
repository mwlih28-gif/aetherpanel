import React, { useState, useEffect } from 'react'
import { HardDrive, Plus, MapPin, Cpu, MemoryStick, Wifi, Settings, Trash2 } from 'lucide-react'

interface Node {
  id: string
  name: string
  hostname: string
  ip: string
  location: string
  status: 'online' | 'offline' | 'maintenance'
  resources: {
    cpu: { used: number; total: number }
    memory: { used: number; total: number }
    disk: { used: number; total: number }
  }
  servers: number
  lastSeen: string
}

export default function NodesPage() {
  const [nodes, setNodes] = useState<Node[]>([])
  const [showAddModal, setShowAddModal] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchNodes()
  }, [])

  const fetchNodes = async () => {
    try {
      // TODO: Replace with actual API call
      setLoading(false)
      setNodes([]) // Start with empty nodes
    } catch (error) {
      console.error('Failed to fetch nodes:', error)
      setLoading(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online': return 'bg-green-100 text-green-800'
      case 'offline': return 'bg-red-100 text-red-800'
      case 'maintenance': return 'bg-yellow-100 text-yellow-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-foreground">Nodes</h1>
            <p className="text-muted-foreground">Manage your server nodes</p>
          </div>
        </div>
        <div className="flex items-center justify-center py-12">
          <div className="text-muted-foreground">Loading nodes...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Nodes</h1>
          <p className="text-muted-foreground">Manage your server nodes</p>
        </div>
        <button 
          onClick={() => setShowAddModal(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2"
        >
          <Plus className="h-4 w-4" />
          Add Node
        </button>
      </div>

      {nodes.length === 0 ? (
        <div className="text-center py-12">
          <HardDrive className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-semibold mb-2">No nodes found</h3>
          <p className="text-muted-foreground mb-4">
            Add your first node to start managing game servers
          </p>
          <button 
            onClick={() => setShowAddModal(true)}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 mx-auto"
          >
            <Plus className="h-4 w-4" />
            Add Node
          </button>
        </div>
      ) : (
        <div className="grid gap-4">
          {nodes.map((node) => (
            <div key={node.id} className="bg-card border rounded-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <HardDrive className={`h-8 w-8 ${node.status === 'online' ? 'text-green-500' : 'text-red-500'}`} />
                  <div>
                    <h3 className="font-semibold text-lg">{node.name}</h3>
                    <p className="text-sm text-muted-foreground">{node.hostname}</p>
                    <div className="flex items-center gap-4 mt-1">
                      <span className="text-xs text-muted-foreground flex items-center gap-1">
                        <Wifi className="h-3 w-3" />
                        {node.ip}
                      </span>
                      <span className="text-xs text-muted-foreground flex items-center gap-1">
                        <MapPin className="h-3 w-3" />
                        {node.location}
                      </span>
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`px-2 py-1 rounded-full text-xs ${getStatusColor(node.status)}`}>
                    {node.status}
                  </span>
                  <button className="p-2 hover:bg-muted rounded-lg">
                    <Settings className="h-4 w-4" />
                  </button>
                  <button className="p-2 hover:bg-muted rounded-lg text-red-500">
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-4 gap-4">
                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="flex items-center gap-2 mb-2">
                    <Cpu className="h-4 w-4 text-blue-500" />
                    <span className="text-sm font-medium">CPU</span>
                  </div>
                  <div className="text-lg font-semibold">
                    {node.resources.cpu.used}%
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div 
                      className="bg-blue-500 h-2 rounded-full" 
                      style={{ width: `${node.resources.cpu.used}%` }}
                    ></div>
                  </div>
                </div>

                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="flex items-center gap-2 mb-2">
                    <MemoryStick className="h-4 w-4 text-green-500" />
                    <span className="text-sm font-medium">Memory</span>
                  </div>
                  <div className="text-lg font-semibold">
                    {Math.round(node.resources.memory.used / 1024)}GB / {Math.round(node.resources.memory.total / 1024)}GB
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div 
                      className="bg-green-500 h-2 rounded-full" 
                      style={{ width: `${(node.resources.memory.used / node.resources.memory.total) * 100}%` }}
                    ></div>
                  </div>
                </div>

                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="flex items-center gap-2 mb-2">
                    <HardDrive className="h-4 w-4 text-purple-500" />
                    <span className="text-sm font-medium">Disk</span>
                  </div>
                  <div className="text-lg font-semibold">
                    {Math.round(node.resources.disk.used / 1024)}GB / {Math.round(node.resources.disk.total / 1024)}GB
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                    <div 
                      className="bg-purple-500 h-2 rounded-full" 
                      style={{ width: `${(node.resources.disk.used / node.resources.disk.total) * 100}%` }}
                    ></div>
                  </div>
                </div>

                <div className="bg-muted/50 rounded-lg p-3">
                  <div className="flex items-center gap-2 mb-2">
                    <HardDrive className="h-4 w-4 text-orange-500" />
                    <span className="text-sm font-medium">Servers</span>
                  </div>
                  <div className="text-lg font-semibold">
                    {node.servers}
                  </div>
                  <div className="text-xs text-muted-foreground mt-1">
                    Active servers
                  </div>
                </div>
              </div>

              <div className="mt-4 pt-4 border-t">
                <span className="text-xs text-muted-foreground">
                  Last seen: {node.lastSeen}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}

      {showAddModal && <AddNodeModal onClose={() => setShowAddModal(false)} onAdd={fetchNodes} />}
    </div>
  )
}

function AddNodeModal({ onClose, onAdd }: { onClose: () => void; onAdd: () => void }) {
  const [formData, setFormData] = useState({
    name: '',
    hostname: '',
    ip: '',
    location: '',
    sshPort: '22',
    sshUser: 'root',
    sshKey: ''
  })
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    
    try {
      // TODO: Replace with actual API call
      console.log('Adding node:', formData)
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      onAdd()
      onClose()
    } catch (error) {
      console.error('Failed to add node:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-card border rounded-lg p-6 w-full max-w-md">
        <h2 className="text-xl font-semibold mb-4">Add New Node</h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Node Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              placeholder="e.g., US-East-1"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Hostname/Domain</label>
            <input
              type="text"
              value={formData.hostname}
              onChange={(e) => setFormData({ ...formData, hostname: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              placeholder="e.g., node1.example.com"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">IP Address</label>
            <input
              type="text"
              value={formData.ip}
              onChange={(e) => setFormData({ ...formData, ip: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              placeholder="e.g., 192.168.1.100"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Location</label>
            <input
              type="text"
              value={formData.location}
              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg"
              placeholder="e.g., New York, USA"
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">SSH Port</label>
              <input
                type="number"
                value={formData.sshPort}
                onChange={(e) => setFormData({ ...formData, sshPort: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg"
                placeholder="22"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">SSH User</label>
              <input
                type="text"
                value={formData.sshUser}
                onChange={(e) => setFormData({ ...formData, sshUser: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg"
                placeholder="root"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">SSH Private Key</label>
            <textarea
              value={formData.sshKey}
              onChange={(e) => setFormData({ ...formData, sshKey: e.target.value })}
              className="w-full px-3 py-2 border rounded-lg h-24"
              placeholder="-----BEGIN PRIVATE KEY-----"
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
              disabled={loading}
              className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg disabled:opacity-50"
            >
              {loading ? 'Adding...' : 'Add Node'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
