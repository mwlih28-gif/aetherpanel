# 🔥 AETHER PANEL 🔥

**Next-Generation Game Server Management Platform**

Aether Panel, Pterodactyl, Multicraft, AMP ve diğer alternatifleri geride bırakan, production-ready, ölçeklenebilir ve modern bir oyun sunucu yönetim panelidir.

---

## 🎯 Özellikler

### Core Features
- **Multi-Game Support**: Minecraft (Java/Bedrock), Rust, ARK, CS2, Valheim, ve 50+ oyun
- **Real-time Console**: WebSocket tabanlı canlı sunucu konsolu
- **Plugin/Mod Marketplace**: CurseForge, Modrinth, Spigot entegrasyonu
- **Player Analytics**: Envanter görüntüleme, chat/komut geçmişi, death logs
- **Automated Backups**: Scheduled ve on-demand backup sistemi
- **Resource Management**: CPU, RAM, Disk, Network limitleri

### Enterprise Features
- **Multi-Node Architecture**: Distributed node sistemi ile sınırsız ölçekleme
- **Reseller System**: White-label reseller paneli
- **Billing Integration**: Kredi sistemi, paketler, otomatik yenileme
- **RBAC**: Role-Based Access Control
- **2FA/MFA**: TOTP tabanlı iki faktörlü doğrulama
- **Audit Logging**: Tüm işlemlerin detaylı kaydı

---

## 🏗️ Mimari

```
┌─────────────────────────────────────────────────────────────────┐
│                        AETHER PANEL                              │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Admin Panel │  │ User Panel  │  │Reseller Panel│             │
│  │   (React)   │  │   (React)   │  │   (React)    │             │
│  └──────┬──────┘  └──────┬──────┘  └──────┬───────┘             │
│         │                │                │                      │
│         └────────────────┼────────────────┘                      │
│                          ▼                                       │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    API GATEWAY                             │  │
│  │              (Traefik / Caddy + Rate Limiting)             │  │
│  └───────────────────────────────────────────────────────────┘  │
│                          │                                       │
│         ┌────────────────┼────────────────┐                      │
│         ▼                ▼                ▼                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Auth Service│  │Server Service│ │Billing Service│            │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Node Service│  │Plugin Service│ │Metrics Service│            │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│                          │                                       │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    MESSAGE BROKER                          │  │
│  │                  (Redis Pub/Sub + Streams)                 │  │
│  └───────────────────────────────────────────────────────────┘  │
│                          │                                       │
│         ┌────────────────┼────────────────┐                      │
│         ▼                ▼                ▼                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ PostgreSQL  │  │    Redis    │  │ TimescaleDB │              │
│  │  (Primary)  │  │(Cache/Queue)│  │  (Metrics)  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
   ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
   │ Node Agent  │ │ Node Agent  │ │ Node Agent  │
   │   (Go)      │ │   (Go)      │ │   (Go)      │
   ├─────────────┤ ├─────────────┤ ├─────────────┤
   │  Docker API │ │  Docker API │ │  Docker API │
   │  Containers │ │  Containers │ │  Containers │
   └─────────────┘ └─────────────┘ └─────────────┘
```

---

## 🚀 Hızlı Kurulum

### Tek Komut Kurulum
```bash
bash <(curl -s https://get.aetherpanel.io/install.sh)
```

### Manuel Kurulum
```bash
git clone https://github.com/aetherpanel/aether-panel.git
cd aether-panel
cp .env.example .env
# .env dosyasını düzenleyin
docker-compose up -d
```

---

## 📁 Proje Yapısı

```
aether-panel/
├── backend/                 # Go Backend (Clean Architecture)
│   ├── cmd/                 # Entry points
│   │   ├── api/            # API server
│   │   ├── agent/          # Node agent
│   │   └── migrate/        # Database migrations
│   ├── internal/           # Private application code
│   │   ├── domain/         # Domain entities & interfaces
│   │   ├── application/    # Use cases & services
│   │   ├── infrastructure/ # External implementations
│   │   └── interfaces/     # API handlers & middleware
│   ├── pkg/                # Public packages
│   └── configs/            # Configuration files
├── frontend/               # React + TypeScript Frontend
│   ├── apps/
│   │   ├── admin/         # Admin panel
│   │   ├── user/          # User panel
│   │   └── reseller/      # Reseller panel
│   └── packages/          # Shared packages
├── agent/                  # Node Agent (standalone)
├── installer/              # Installation scripts
├── docker/                 # Docker configurations
├── docs/                   # Documentation
└── deploy/                 # Deployment configs
```

---

## 🛠️ Teknoloji Stack

### Backend
- **Language**: Go 1.22+
- **Framework**: Fiber v2 (High-performance)
- **ORM**: GORM + sqlx
- **Validation**: go-playground/validator
- **Auth**: JWT + PASETO

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **Styling**: TailwindCSS + shadcn/ui
- **State**: Zustand
- **Real-time**: Socket.io

### Database
- **Primary**: PostgreSQL 16
- **Cache**: Redis 7
- **Metrics**: TimescaleDB

### Infrastructure
- **Container**: Docker + Docker Compose
- **Reverse Proxy**: Traefik / Caddy
- **SSL**: Let's Encrypt
- **Monitoring**: Prometheus + Grafana

---

## 📖 Dokümantasyon

- [Kurulum Rehberi](docs/installation.md)
- [API Referansı](docs/api-reference.md)
- [Node Agent](docs/node-agent.md)
- [Plugin Geliştirme](docs/plugin-development.md)
- [Güvenlik](docs/security.md)

---

## 🔒 Güvenlik

- JWT + Refresh Token rotasyonu
- TOTP 2FA desteği
- Rate limiting (per-IP, per-user)
- IP whitelist/blacklist
- Encrypted secrets (AES-256-GCM)
- Full audit logging
- RBAC permission system

---

## 📊 Desteklenen Oyunlar

| Oyun | Java | Bedrock | Modded |
|------|------|---------|--------|
| Minecraft | ✅ | ✅ | ✅ |
| Rust | ✅ | - | ✅ |
| ARK: Survival | ✅ | - | ✅ |
| CS2 | ✅ | - | - |
| Valheim | ✅ | - | ✅ |
| Terraria | ✅ | - | ✅ |
| 7 Days to Die | ✅ | - | ✅ |

---

## 📄 Lisans

Aether Panel, [MIT License](LICENSE) altında lisanslanmıştır.

---

## 🤝 Katkıda Bulunma

Katkılarınızı bekliyoruz! Lütfen [CONTRIBUTING.md](CONTRIBUTING.md) dosyasını inceleyin.

---

**Made with ❤️ by Aether Team**
