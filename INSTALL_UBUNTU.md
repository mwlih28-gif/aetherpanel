# 🔥 Aether Panel - Ubuntu VPS Kurulum Rehberi

## Gereksinimler

- **OS**: Ubuntu 22.04 LTS veya 24.04 LTS
- **RAM**: Minimum 2GB (4GB önerilir)
- **CPU**: 2 vCPU
- **Disk**: 20GB SSD
- **Ağ**: Açık portlar: 80, 443, 8080, 3000

---

## Hızlı Kurulum (Tek Komut)

```bash
curl -sSL https://raw.githubusercontent.com/aetherpanel/aether-panel/main/install.sh | sudo bash
```

---

## Manuel Kurulum

### 1. Sistemi Güncelle

```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Gerekli Paketleri Kur

```bash
sudo apt install -y curl wget git ca-certificates gnupg lsb-release openssl
```

### 3. Docker Kurulumu

```bash
# Docker GPG key ekle
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Docker repository ekle
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Docker kur
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Docker'ı başlat
sudo systemctl enable docker
sudo systemctl start docker

# Kullanıcıyı docker grubuna ekle (opsiyonel)
sudo usermod -aG docker $USER
```

### 4. Aether Panel'i İndir

```bash
# Kurulum dizini oluştur
sudo mkdir -p /opt/aether-panel
cd /opt/aether-panel

# Dosyaları kopyala (git clone veya scp ile)
# git clone https://github.com/aetherpanel/aether-panel.git .
# VEYA dosyaları manuel olarak yükle
```

### 5. Environment Dosyasını Oluştur

```bash
# Güvenli şifreler oluştur
DB_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)
REDIS_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)
JWT_SECRET=$(openssl rand -base64 64 | tr -dc 'a-zA-Z0-9' | head -c 64)
ENCRYPTION_KEY=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)

# .env dosyası oluştur
cat > /opt/aether-panel/.env << EOF
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

# SSL (opsiyonel)
ACME_EMAIL=admin@example.com
EOF

# Dosya izinlerini ayarla
sudo chmod 600 /opt/aether-panel/.env
```

### 6. Veri Dizinlerini Oluştur

```bash
sudo mkdir -p /opt/aether-panel/data/{backups,logs,servers}
sudo mkdir -p /var/lib/aether/{backups,servers}
```

### 7. Panel'i Başlat

```bash
cd /opt/aether-panel
sudo docker compose up -d
```

### 8. Kurulumu Doğrula

```bash
# Container'ları kontrol et
sudo docker compose ps

# Logları izle
sudo docker compose logs -f

# API sağlık kontrolü
curl http://localhost:8080/health
```

---

## Erişim

Panel kurulduktan sonra:

- **Panel URL**: `http://SUNUCU_IP:3000`
- **API URL**: `http://SUNUCU_IP:8080`

---

## SSL Sertifikası (Let's Encrypt)

### Traefik ile Otomatik SSL

```bash
# .env dosyasına domain ekle
echo "DOMAIN=panel.example.com" >> /opt/aether-panel/.env
echo "ACME_EMAIL=admin@example.com" >> /opt/aether-panel/.env

# Traefik profili ile başlat
cd /opt/aether-panel
sudo docker compose --profile with-traefik up -d
```

### Nginx ile Manuel SSL

```bash
# Certbot kur
sudo apt install -y certbot

# Sertifika al
sudo certbot certonly --standalone -d panel.example.com

# Nginx kur ve yapılandır
sudo apt install -y nginx
```

---

## Yönetim Komutları

```bash
# Servisleri durdur
cd /opt/aether-panel && sudo docker compose down

# Servisleri yeniden başlat
cd /opt/aether-panel && sudo docker compose restart

# Logları görüntüle
cd /opt/aether-panel && sudo docker compose logs -f api

# Güncelleme
cd /opt/aether-panel
sudo docker compose pull
sudo docker compose up -d

# Veritabanı yedeği
sudo docker exec aether_postgres pg_dump -U aether aether_panel > backup.sql
```

---

## Firewall Ayarları

```bash
# UFW ile port aç
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 3000/tcp  # Frontend (geçici)
sudo ufw allow 8080/tcp  # API (geçici)
sudo ufw enable
```

---

## Sorun Giderme

### Container başlamıyor
```bash
sudo docker compose logs api
sudo docker compose logs postgres
```

### Veritabanı bağlantı hatası
```bash
# PostgreSQL'in hazır olmasını bekle
sudo docker compose restart api
```

### Port çakışması
```bash
# Kullanılan portları kontrol et
sudo netstat -tlnp | grep -E '(3000|8080)'
```

---

## Dosya Yapısı

```
/opt/aether-panel/
├── .env                 # Ortam değişkenleri
├── docker-compose.yml   # Docker yapılandırması
├── backend/             # Go API
├── frontend/            # React UI
├── agent/               # Node Agent
└── data/
    ├── backups/         # Yedekler
    ├── logs/            # Loglar
    └── servers/         # Sunucu verileri
```

---

## Destek

Sorun yaşarsanız:
1. Logları kontrol edin: `docker compose logs`
2. GitHub Issues açın
3. Dokümantasyonu inceleyin
