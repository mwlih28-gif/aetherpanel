#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Variables
INSTALL_DIR="/opt/aether-panel"
LOG_FILE="/var/log/aether-install.log"
COMPOSE_VERSION="2.23.0"

# Functions
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] INFO: $1" >> "$LOG_FILE"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] WARN: $1" >> "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] ERROR: $1" >> "$LOG_FILE"
    exit 1
}

banner() {
    echo -e "${BLUE}"
    echo "╔═══════════════════════════════════════════════════════════╗"
    echo "║                                                           ║"
    echo "║     🔥 AETHER PANEL INSTALLER 🔥                         ║"
    echo "║     Next-Generation Game Server Management               ║"
    echo "║                                                           ║"
    echo "╚═══════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root"
    fi
}

check_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$ID
        VERSION=$VERSION_ID
    else
        error "Cannot detect OS"
    fi

    case $OS in
        ubuntu)
            if [[ "$VERSION" != "22.04" && "$VERSION" != "24.04" ]]; then
                warn "Ubuntu $VERSION is not officially supported. Recommended: 22.04 LTS"
            fi
            ;;
        debian)
            if [[ "$VERSION" != "11" && "$VERSION" != "12" ]]; then
                warn "Debian $VERSION is not officially supported. Recommended: 11 or 12"
            fi
            ;;
        *)
            error "Unsupported OS: $OS. Supported: Ubuntu 22.04+, Debian 11+"
            ;;
    esac

    log "Detected OS: $OS $VERSION"
}

install_dependencies() {
    log "Installing dependencies..."
    
    apt-get update -qq
    apt-get install -y -qq \
        curl \
        wget \
        git \
        ca-certificates \
        gnupg \
        lsb-release \
        openssl \
        jq
}

install_docker() {
    if command -v docker &> /dev/null; then
        log "Docker is already installed"
        return
    fi

    log "Installing Docker..."
    
    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/$OS/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    # Add the repository
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$OS \
        $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

    apt-get update -qq
    apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    systemctl enable docker
    systemctl start docker

    log "Docker installed successfully"
}

generate_secrets() {
    log "Generating secrets..."
    
    DB_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)
    REDIS_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)
    JWT_SECRET=$(openssl rand -base64 64 | tr -dc 'a-zA-Z0-9' | head -c 64)
    ENCRYPTION_KEY=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)
}

setup_panel() {
    log "Setting up Aether Panel..."

    # Create installation directory
    mkdir -p "$INSTALL_DIR"
    cd "$INSTALL_DIR"

    # Clone or copy files
    if [[ -d "/tmp/aether-panel" ]]; then
        cp -r /tmp/aether-panel/* "$INSTALL_DIR/"
    else
        log "Downloading Aether Panel..."
        # In production, this would clone from git
        # git clone https://github.com/aetherpanel/aether-panel.git .
    fi

    # Create .env file
    cat > "$INSTALL_DIR/.env" << EOF
# Aether Panel Configuration
# Generated on $(date)

# Database
DB_USER=aether
DB_PASSWORD=$DB_PASSWORD
DB_NAME=aether_panel

# Redis
REDIS_PASSWORD=$REDIS_PASSWORD

# Security
JWT_SECRET=$JWT_SECRET
ENCRYPTION_KEY=$ENCRYPTION_KEY

# Ports
API_PORT=8080
FRONTEND_PORT=3000

# SSL (optional)
ACME_EMAIL=admin@example.com
EOF

    chmod 600 "$INSTALL_DIR/.env"

    # Create data directories
    mkdir -p "$INSTALL_DIR/data/backups"
    mkdir -p "$INSTALL_DIR/data/logs"
    mkdir -p "$INSTALL_DIR/data/servers"

    log "Panel files configured"
}

start_services() {
    log "Starting Aether Panel services..."
    
    cd "$INSTALL_DIR"
    docker compose up -d

    log "Waiting for services to start..."
    sleep 10

    # Check if services are running
    if docker compose ps | grep -q "Up"; then
        log "Services started successfully"
    else
        error "Failed to start services. Check logs with: docker compose logs"
    fi
}

create_admin_user() {
    log "Creating admin user..."
    
    read -p "Enter admin email: " ADMIN_EMAIL
    read -s -p "Enter admin password: " ADMIN_PASSWORD
    echo

    # In production, this would call the API to create the admin user
    log "Admin user will be created on first login"
}

print_summary() {
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║          🎉 INSTALLATION COMPLETE! 🎉                     ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "Panel URL:     ${BLUE}http://$(hostname -I | awk '{print $1}'):3000${NC}"
    echo -e "API URL:       ${BLUE}http://$(hostname -I | awk '{print $1}'):8080${NC}"
    echo ""
    echo -e "Installation directory: ${YELLOW}$INSTALL_DIR${NC}"
    echo -e "Log file: ${YELLOW}$LOG_FILE${NC}"
    echo ""
    echo -e "${YELLOW}Important:${NC}"
    echo "  - Credentials are stored in: $INSTALL_DIR/.env"
    echo "  - To view logs: cd $INSTALL_DIR && docker compose logs -f"
    echo "  - To stop: cd $INSTALL_DIR && docker compose down"
    echo "  - To update: cd $INSTALL_DIR && docker compose pull && docker compose up -d"
    echo ""
}

# Main installation flow
main() {
    banner
    
    # Create log file
    mkdir -p "$(dirname "$LOG_FILE")"
    touch "$LOG_FILE"
    
    log "Starting Aether Panel installation..."
    
    check_root
    check_os
    install_dependencies
    install_docker
    generate_secrets
    setup_panel
    start_services
    
    print_summary
    
    log "Installation completed successfully!"
}

# Run main function
main "$@"
