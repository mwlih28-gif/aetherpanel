import React, { useState, useEffect } from 'react'
import { HardDrive, Plus, MapPin, Cpu, MemoryStick, Wifi, Settings, Trash2, Key, Copy, Eye, EyeOff, Terminal, Globe } from 'lucide-react'

interface Location {
  id: string
  short: string
  long: string
}

interface Node {
  id: string
  uuid: string
  name: string
  description: string
  locationId: string
  location?: Location
  fqdn: string
  scheme: 'http' | 'https'
  behindProxy: boolean
  maintenanceMode: boolean
  memory: number
  memoryOverallocate: number
  disk: number
  diskOverallocate: number
  uploadSize: number
  daemonListenPort: number
  daemonSftpPort: number
  daemonToken: string
  publicKey: boolean
  status: 'online' | 'offline' | 'maintenance'
  resources: {
    cpu: { used: number; total: number }
    memory: { used: number; total: number }
    disk: { used: number; total: number }
  }
  servers: number
  lastSeen: string
  createdAt: string
  updatedAt: string
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
  const [currentStep, setCurrentStep] = useState(1)
  const [locations, setLocations] = useState<Location[]>([])
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    locationId: '',
    fqdn: '',
    scheme: 'https' as 'http' | 'https',
    behindProxy: false,
    publicKey: true,
    memory: 1024,
    memoryOverallocate: 0,
    disk: 10240,
    diskOverallocate: 0,
    uploadSize: 100,
    daemonListenPort: 8080,
    daemonSftpPort: 2022
  })
  const [generatedToken, setGeneratedToken] = useState('')
  const [showToken, setShowToken] = useState(false)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetchLocations()
  }, [])

  const fetchLocations = async () => {
    try {
      // TODO: Replace with actual API call
      setLocations([])
    } catch (error) {
      console.error('Failed to fetch locations:', error)
    }
  }

  const generateToken = () => {
    const token = Array.from(crypto.getRandomValues(new Uint8Array(32)))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('')
    setGeneratedToken(token)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    
    try {
      if (!generatedToken) {
        generateToken()
      }
      
      // TODO: Replace with actual API call
      console.log('Creating node:', { ...formData, daemonToken: generatedToken })
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setCurrentStep(3) // Show installation instructions
    } catch (error) {
      console.error('Failed to create node:', error)
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

  const getInstallationScript = () => {
    return `#!/bin/bash
# Aether Panel Wings Installation Script
# Automated installation script for Wings daemon

set -e

# Colors for output
RED='\\033[0;31m'
GREEN='\\033[0;32m'
YELLOW='\\033[1;33m'
BLUE='\\033[0;34m'
NC='\\033[0m'

print_info() { echo -e "\${BLUE}[INFO]\${NC} \$1"; }
print_success() { echo -e "\${GREEN}[SUCCESS]\${NC} \$1"; }
print_warning() { echo -e "\${YELLOW}[WARNING]\${NC} \$1"; }
print_error() { echo -e "\${RED}[ERROR]\${NC} \$1"; }

# Configuration
PANEL_URL="${window.location.protocol}//${window.location.host}"
NODE_TOKEN="${generatedToken}"
DAEMON_PORT="${formData.daemonListenPort}"
SFTP_PORT="${formData.daemonSftpPort}"
NODE_FQDN="${formData.fqdn}"
USE_SSL="${formData.scheme === 'https' ? 'true' : 'false'}"

echo -e "\${GREEN}=== Aether Panel Wings Installer ===\${NC}"
echo

# Check if running as root
if [[ \$EUID -ne 0 ]]; then
    print_error "This script must be run as root"
    exit 1
fi

print_info "Installing dependencies..."
if command -v apt &> /dev/null; then
    apt update && apt install -y curl wget tar unzip software-properties-common apt-transport-https ca-certificates gnupg lsb-release
elif command -v yum &> /dev/null; then
    yum update -y && yum install -y curl wget tar unzip yum-utils device-mapper-persistent-data lvm2
else
    print_error "Unsupported package manager"
    exit 1
fi

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    print_info "Installing Docker..."
    if command -v apt &> /dev/null; then
        curl -fsSL https://get.docker.com | sh
    else
        curl -fsSL https://get.docker.com | sh
    fi
    systemctl start docker && systemctl enable docker
    print_success "Docker installed"
else
    print_info "Docker already installed"
fi

# Create pterodactyl user
print_info "Creating pterodactyl user..."
if ! id "pterodactyl" &>/dev/null; then
    useradd -r -d /etc/pterodactyl -s /usr/sbin/nologin pterodactyl
fi

# Create directories
print_info "Creating directories..."
mkdir -p /etc/pterodactyl /var/log/pterodactyl /var/lib/pterodactyl/{volumes,backups} /var/run/wings
chown -R pterodactyl:pterodactyl /etc/pterodactyl /var/log/pterodactyl /var/lib/pterodactyl /var/run/wings

# Download Wings
print_info "Downloading Wings..."
WINGS_VERSION=\$(curl -s https://api.github.com/repos/pterodactyl/wings/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\\1/')
curl -L -o /usr/local/bin/wings "https://github.com/pterodactyl/wings/releases/download/\${WINGS_VERSION}/wings_linux_amd64"
chmod u+x /usr/local/bin/wings
print_success "Wings downloaded: \$WINGS_VERSION"

# Create configuration
print_info "Creating configuration..."
NODE_UUID=\$(cat /proc/sys/kernel/random/uuid)

cat > /etc/pterodactyl/config.yml << 'EOF'
debug: false
uuid: \${NODE_UUID}
token_id: ${generatedToken.substring(0, 16)}
token: ${generatedToken}
api:
  host: ${formData.fqdn}
  port: ${formData.daemonListenPort}
  ssl:
    enabled: ${formData.scheme === 'https'}
    cert: /etc/letsencrypt/live/${formData.fqdn}/fullchain.pem
    key: /etc/letsencrypt/live/${formData.fqdn}/privkey.pem
system:
  data: /var/lib/pterodactyl/volumes
  sftp:
    bind_port: ${formData.daemonSftpPort}
allowed_mounts: []
remote: ${window.location.protocol}//${window.location.host}
EOF

chown pterodactyl:pterodactyl /etc/pterodactyl/config.yml
chmod 600 /etc/pterodactyl/config.yml

# Create systemd service
print_info "Creating systemd service..."
cat > /etc/systemd/system/wings.service << 'EOF'
[Unit]
Description=Pterodactyl Wings Daemon
After=docker.service network-online.target
Requires=docker.service
PartOf=docker.service
DefaultDependencies=no

[Service]
User=pterodactyl
WorkingDirectory=/etc/pterodactyl
LimitNOFILE=4096
PIDFile=/var/run/wings/daemon.pid
ExecStart=/usr/local/bin/wings
Restart=on-failure
StartLimitInterval=180
StartLimitBurst=30
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload

# Configure firewall
print_info "Configuring firewall..."
if command -v ufw &> /dev/null && ufw status | grep -q "Status: active"; then
    ufw allow ${formData.daemonListenPort}/tcp
    ufw allow ${formData.daemonSftpPort}/tcp
elif command -v firewall-cmd &> /dev/null; then
    firewall-cmd --permanent --add-port=${formData.daemonListenPort}/tcp
    firewall-cmd --permanent --add-port=${formData.daemonSftpPort}/tcp
    firewall-cmd --reload
fi

# Start Wings
print_info "Starting Wings service..."
systemctl enable wings
systemctl start wings

sleep 3

if systemctl is-active --quiet wings; then
    print_success "Wings installation completed successfully!"
    echo
    echo -e "\${GREEN}=== Installation Summary ===\${NC}"
    echo -e "\${BLUE}Panel URL:\${NC} ${window.location.protocol}//${window.location.host}"
    echo -e "\${BLUE}Node FQDN:\${NC} ${formData.fqdn}"
    echo -e "\${BLUE}Daemon Port:\${NC} ${formData.daemonListenPort}"
    echo -e "\${BLUE}SFTP Port:\${NC} ${formData.daemonSftpPort}"
    echo
    echo -e "\${YELLOW}Useful Commands:\${NC}"
    echo -e "  View logs: \${BLUE}journalctl -u wings -f\${NC}"
    echo -e "  Restart: \${BLUE}systemctl restart wings\${NC}"
    echo -e "  Status: \${BLUE}systemctl status wings\${NC}"
    echo
    echo -e "\${GREEN}Your node should now appear as online in Aether Panel!\${NC}"
else
    print_error "Failed to start Wings service"
    echo "Check logs with: journalctl -u wings -f"
    exit 1
fi`
  }

  if (currentStep === 3) {
    return (
      <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
        <div className="bg-card border rounded-lg p-6 w-full max-w-4xl max-h-[90vh] overflow-y-auto">
          <h2 className="text-xl font-semibold mb-4">Wings Installation</h2>
          
          <div className="space-y-6">
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h3 className="font-semibold text-blue-900 mb-2">Installation Instructions</h3>
              <p className="text-blue-800 text-sm">
                Run the following command on your server to automatically install and configure Wings daemon.
                The script will handle all dependencies including Docker installation.
              </p>
            </div>

            <div>
              <div className="flex items-center justify-between mb-2">
                <h4 className="font-medium">Quick Install (Recommended)</h4>
                <button
                  onClick={() => copyToClipboard(`bash <(curl -s ${window.location.protocol}//${window.location.host}/install-wings.sh) --token="${generatedToken}" --fqdn="${formData.fqdn}" --daemon-port="${formData.daemonListenPort}" --sftp-port="${formData.daemonSftpPort}" --ssl="${formData.scheme === 'https'}"`)}
                  className="flex items-center gap-2 px-3 py-1 bg-green-100 hover:bg-green-200 rounded text-sm"
                >
                  <Copy className="h-4 w-4" />
                  Copy Command
                </button>
              </div>
              <div className="bg-gray-900 text-green-400 p-4 rounded-lg text-sm font-mono">
                bash &lt;(curl -s {window.location.protocol}//{window.location.host}/install-wings.sh) --token="{generatedToken}" --fqdn="{formData.fqdn}" --daemon-port="{formData.daemonListenPort}" --sftp-port="{formData.daemonSftpPort}" --ssl="{formData.scheme === 'https'}"
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                This command will automatically download and run the installation script with your node configuration.
              </p>
            </div>

            <div>
              <div className="flex items-center justify-between mb-2">
                <h4 className="font-medium">Installation Script</h4>
                <button
                  onClick={() => copyToClipboard(getInstallationScript())}
                  className="flex items-center gap-2 px-3 py-1 bg-gray-100 hover:bg-gray-200 rounded text-sm"
                >
                  <Copy className="h-4 w-4" />
                  Copy Script
                </button>
              </div>
              <pre className="bg-gray-900 text-green-400 p-4 rounded-lg text-sm overflow-x-auto">
                {getInstallationScript()}
              </pre>
            </div>

            <div>
              <h4 className="font-medium mb-2">Node Token</h4>
              <div className="flex items-center gap-2">
                <input
                  type={showToken ? 'text' : 'password'}
                  value={generatedToken}
                  readOnly
                  className="flex-1 px-3 py-2 border rounded-lg font-mono text-sm"
                />
                <button
                  onClick={() => setShowToken(!showToken)}
                  className="p-2 border rounded-lg hover:bg-muted"
                >
                  {showToken ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
                <button
                  onClick={() => copyToClipboard(generatedToken)}
                  className="p-2 border rounded-lg hover:bg-muted"
                >
                  <Copy className="h-4 w-4" />
                </button>
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                This token is used to authenticate the Wings daemon with the panel.
              </p>
            </div>

            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <h4 className="font-semibold text-yellow-900 mb-2">Important Notes</h4>
              <ul className="text-yellow-800 text-sm space-y-1">
                <li>• Ensure Docker is installed and running on the target server</li>
                <li>• Make sure ports {formData.daemonListenPort} and {formData.daemonSftpPort} are open</li>
                <li>• The server should be able to reach this panel at {window.location.host}</li>
                <li>• SSL certificates are required if using HTTPS scheme</li>
              </ul>
            </div>

            <div className="flex gap-3 pt-4">
              <button
                onClick={() => {
                  onAdd()
                  onClose()
                }}
                className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg"
              >
                Complete Setup
              </button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-card border rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <h2 className="text-xl font-semibold mb-4">Create New Node</h2>
        
        {/* Step Indicator */}
        <div className="flex items-center mb-6">
          <div className={`flex items-center justify-center w-8 h-8 rounded-full ${currentStep >= 1 ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}>
            1
          </div>
          <div className={`flex-1 h-1 mx-2 ${currentStep >= 2 ? 'bg-blue-500' : 'bg-gray-200'}`}></div>
          <div className={`flex items-center justify-center w-8 h-8 rounded-full ${currentStep >= 2 ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}>
            2
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Basic Details */}
          <div>
            <h3 className="font-semibold mb-4">Basic Details</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Name</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg"
                  placeholder="Node Name"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Location</label>
                <select
                  value={formData.locationId}
                  onChange={(e) => setFormData({ ...formData, locationId: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg"
                  required
                >
                  <option value="">Select a Location</option>
                  {locations.map(location => (
                    <option key={location.id} value={location.id}>
                      {location.short} - {location.long}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div className="mt-4">
              <label className="block text-sm font-medium mb-1">Description</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="w-full px-3 py-2 border rounded-lg h-20"
                placeholder="A brief description of this server node"
              />
            </div>
          </div>

          {/* Configuration */}
          <div>
            <h3 className="font-semibold mb-4">Configuration</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">FQDN</label>
                <input
                  type="text"
                  value={formData.fqdn}
                  onChange={(e) => setFormData({ ...formData, fqdn: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg"
                  placeholder="node.example.com"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Communicate Over SSL</label>
                <div className="flex gap-4 mt-2">
                  <label className="flex items-center">
                    <input
                      type="radio"
                      name="scheme"
                      value="https"
                      checked={formData.scheme === 'https'}
                      onChange={(e) => setFormData({ ...formData, scheme: e.target.value as 'https' })}
                      className="mr-2"
                    />
                    Use SSL Connection
                  </label>
                  <label className="flex items-center">
                    <input
                      type="radio"
                      name="scheme"
                      value="http"
                      checked={formData.scheme === 'http'}
                      onChange={(e) => setFormData({ ...formData, scheme: e.target.value as 'http' })}
                      className="mr-2"
                    />
                    Use HTTP Connection
                  </label>
                </div>
              </div>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
              <div>
                <label className="block text-sm font-medium mb-1">Daemon Port</label>
                <input
                  type="number"
                  value={formData.daemonListenPort}
                  onChange={(e) => setFormData({ ...formData, daemonListenPort: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border rounded-lg"
                  min="1024"
                  max="65535"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Daemon SFTP Port</label>
                <input
                  type="number"
                  value={formData.daemonSftpPort}
                  onChange={(e) => setFormData({ ...formData, daemonSftpPort: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border rounded-lg"
                  min="1024"
                  max="65535"
                />
              </div>
            </div>

            <div className="mt-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={formData.behindProxy}
                  onChange={(e) => setFormData({ ...formData, behindProxy: e.target.checked })}
                  className="mr-2"
                />
                Behind Proxy
              </label>
              <p className="text-xs text-muted-foreground mt-1">
                If you are running the daemon behind a proxy such as Cloudflare, select this to have the daemon skip looking for certificates on boot.
              </p>
            </div>
          </div>

          {/* Resource Limits */}
          <div>
            <h3 className="font-semibold mb-4">Resource Management</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Total Memory (MB)</label>
                <input
                  type="number"
                  value={formData.memory}
                  onChange={(e) => setFormData({ ...formData, memory: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border rounded-lg"
                  min="128"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Memory Over-Allocation (%)</label>
                <input
                  type="number"
                  value={formData.memoryOverallocate}
                  onChange={(e) => setFormData({ ...formData, memoryOverallocate: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border rounded-lg"
                  min="0"
                  max="500"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
              <div>
                <label className="block text-sm font-medium mb-1">Disk Space (MB)</label>
                <input
                  type="number"
                  value={formData.disk}
                  onChange={(e) => setFormData({ ...formData, disk: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border rounded-lg"
                  min="1024"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Disk Over-Allocation (%)</label>
                <input
                  type="number"
                  value={formData.diskOverallocate}
                  onChange={(e) => setFormData({ ...formData, diskOverallocate: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border rounded-lg"
                  min="0"
                  max="500"
                />
              </div>
            </div>
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
              {loading ? 'Creating...' : 'Create Node'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
