#!/bin/bash

# Aether Panel Wings Installation Script
# This script installs and configures Wings daemon for Aether Panel
# Usage: curl -sSL https://raw.githubusercontent.com/your-repo/aether-panel/main/scripts/install-wings.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration variables (will be replaced by panel)
PANEL_URL="${PANEL_URL:-http://localhost:3001}"
NODE_TOKEN="${NODE_TOKEN:-}"
DAEMON_PORT="${DAEMON_PORT:-8080}"
SFTP_PORT="${SFTP_PORT:-2022}"
NODE_FQDN="${NODE_FQDN:-}"
USE_SSL="${USE_SSL:-false}"

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root"
        exit 1
    fi
}

check_os() {
    if [[ ! -f /etc/os-release ]]; then
        print_error "Cannot determine OS. /etc/os-release not found."
        exit 1
    fi
    
    . /etc/os-release
    
    case "$ID" in
        ubuntu|debian)
            PACKAGE_MANAGER="apt"
            ;;
        centos|rhel|fedora)
            PACKAGE_MANAGER="yum"
            ;;
        *)
            print_error "Unsupported OS: $ID"
            exit 1
            ;;
    esac
    
    print_info "Detected OS: $PRETTY_NAME"
}

install_dependencies() {
    print_info "Installing dependencies..."
    
    case "$PACKAGE_MANAGER" in
        apt)
            apt update
            apt install -y curl wget tar unzip software-properties-common apt-transport-https ca-certificates gnupg lsb-release
            ;;
        yum)
            yum update -y
            yum install -y curl wget tar unzip yum-utils device-mapper-persistent-data lvm2
            ;;
    esac
    
    print_success "Dependencies installed"
}

install_docker() {
    if command -v docker &> /dev/null; then
        print_info "Docker is already installed"
        return
    fi
    
    print_info "Installing Docker..."
    
    case "$PACKAGE_MANAGER" in
        apt)
            # Add Docker's official GPG key
            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
            
            # Add Docker repository
            echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
            
            # Install Docker
            apt update
            apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
            ;;
        yum)
            # Add Docker repository
            yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
            
            # Install Docker
            yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
            ;;
    esac
    
    # Start and enable Docker
    systemctl start docker
    systemctl enable docker
    
    print_success "Docker installed and started"
}

create_user() {
    print_info "Creating pterodactyl user..."
    
    if id "pterodactyl" &>/dev/null; then
        print_info "User pterodactyl already exists"
    else
        useradd -r -d /etc/pterodactyl -s /usr/sbin/nologin pterodactyl
        print_success "User pterodactyl created"
    fi
}

create_directories() {
    print_info "Creating directories..."
    
    mkdir -p /etc/pterodactyl
    mkdir -p /var/log/pterodactyl
    mkdir -p /var/lib/pterodactyl/volumes
    mkdir -p /var/lib/pterodactyl/backups
    mkdir -p /var/run/wings
    
    chown -R pterodactyl:pterodactyl /etc/pterodactyl /var/log/pterodactyl /var/lib/pterodactyl
    chown pterodactyl:pterodactyl /var/run/wings
    
    print_success "Directories created"
}

download_wings() {
    print_info "Downloading Wings..."
    
    # Get latest Wings version
    WINGS_VERSION=$(curl -s https://api.github.com/repos/pterodactyl/wings/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [[ -z "$WINGS_VERSION" ]]; then
        print_error "Failed to get Wings version"
        exit 1
    fi
    
    print_info "Latest Wings version: $WINGS_VERSION"
    
    # Download Wings
    curl -L -o /usr/local/bin/wings "https://github.com/pterodactyl/wings/releases/download/${WINGS_VERSION}/wings_linux_amd64"
    chmod u+x /usr/local/bin/wings
    
    print_success "Wings downloaded and installed"
}

create_config() {
    print_info "Creating Wings configuration..."
    
    if [[ -z "$NODE_TOKEN" ]]; then
        print_error "NODE_TOKEN is required but not provided"
        exit 1
    fi
    
    if [[ -z "$NODE_FQDN" ]]; then
        print_error "NODE_FQDN is required but not provided"
        exit 1
    fi
    
    # Generate UUID for this node
    NODE_UUID=$(cat /proc/sys/kernel/random/uuid)
    
    # Create configuration file
    cat > /etc/pterodactyl/config.yml << EOF
debug: false
uuid: ${NODE_UUID}
token_id: ${NODE_TOKEN:0:16}
token: ${NODE_TOKEN}
api:
  host: ${NODE_FQDN}
  port: ${DAEMON_PORT}
  ssl:
    enabled: ${USE_SSL}
    cert: /etc/letsencrypt/live/${NODE_FQDN}/fullchain.pem
    key: /etc/letsencrypt/live/${NODE_FQDN}/privkey.pem
system:
  data: /var/lib/pterodactyl/volumes
  sftp:
    bind_port: ${SFTP_PORT}
allowed_mounts: []
remote: ${PANEL_URL}
EOF
    
    chown pterodactyl:pterodactyl /etc/pterodactyl/config.yml
    chmod 600 /etc/pterodactyl/config.yml
    
    print_success "Configuration file created"
}

create_systemd_service() {
    print_info "Creating systemd service..."
    
    cat > /etc/systemd/system/wings.service << EOF
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
    
    print_success "Systemd service created"
}

configure_firewall() {
    print_info "Configuring firewall..."
    
    # Check if ufw is installed and active
    if command -v ufw &> /dev/null && ufw status | grep -q "Status: active"; then
        print_info "Configuring UFW firewall..."
        ufw allow ${DAEMON_PORT}/tcp
        ufw allow ${SFTP_PORT}/tcp
        ufw allow 2376/tcp  # Docker daemon
        print_success "UFW firewall configured"
    elif command -v firewall-cmd &> /dev/null; then
        print_info "Configuring firewalld..."
        firewall-cmd --permanent --add-port=${DAEMON_PORT}/tcp
        firewall-cmd --permanent --add-port=${SFTP_PORT}/tcp
        firewall-cmd --permanent --add-port=2376/tcp
        firewall-cmd --reload
        print_success "Firewalld configured"
    else
        print_warning "No supported firewall found. Please manually open ports ${DAEMON_PORT} and ${SFTP_PORT}"
    fi
}

start_wings() {
    print_info "Starting Wings service..."
    
    systemctl enable wings
    systemctl start wings
    
    # Wait a moment for service to start
    sleep 3
    
    if systemctl is-active --quiet wings; then
        print_success "Wings service started successfully"
    else
        print_error "Failed to start Wings service"
        print_info "Check logs with: journalctl -u wings -f"
        exit 1
    fi
}

show_status() {
    print_info "Installation completed!"
    echo
    echo -e "${GREEN}=== Aether Panel Wings Installation Summary ===${NC}"
    echo -e "${BLUE}Panel URL:${NC} $PANEL_URL"
    echo -e "${BLUE}Node FQDN:${NC} $NODE_FQDN"
    echo -e "${BLUE}Daemon Port:${NC} $DAEMON_PORT"
    echo -e "${BLUE}SFTP Port:${NC} $SFTP_PORT"
    echo -e "${BLUE}SSL Enabled:${NC} $USE_SSL"
    echo
    echo -e "${YELLOW}Useful Commands:${NC}"
    echo -e "  View logs: ${BLUE}journalctl -u wings -f${NC}"
    echo -e "  Restart service: ${BLUE}systemctl restart wings${NC}"
    echo -e "  Check status: ${BLUE}systemctl status wings${NC}"
    echo -e "  Configuration: ${BLUE}/etc/pterodactyl/config.yml${NC}"
    echo
    echo -e "${GREEN}The node should now appear as online in your Aether Panel!${NC}"
}

main() {
    echo -e "${GREEN}=== Aether Panel Wings Installer ===${NC}"
    echo
    
    check_root
    check_os
    install_dependencies
    install_docker
    create_user
    create_directories
    download_wings
    create_config
    create_systemd_service
    configure_firewall
    start_wings
    show_status
}

# Run main function
main "$@"
