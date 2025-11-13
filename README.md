# 🎮 Tower Defense Game

A modern, full-stack tower defense game built with **Go** (backend) and **React + TypeScript** (frontend). Features real-time gameplay, multiple tower and enemy types, save/load functionality, and a clean ECS-inspired architecture.

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

<img width="1470" height="770" alt="Zrzut ekranu 2025-10-19 o 17 20 48" src="https://github.com/user-attachments/assets/0d492d75-62d9-4cb4-bf67-7b743e238bc0" />




## ✨ Features

### 🎯 Gameplay
- **3 Tower Types**: Basic (balanced), Sniper (long-range), Splash (area damage)
- **4 Enemy Types**: Basic, Fast (2x speed), Tank (high HP), Boss (waves 10, 20, 30...)
- **Dynamic Wave System**: Progressive difficulty with HP/count scaling
- **Save/Load**: Full game state persistence (localStorage + file upload/download)
- **Real-time Updates**: WebSocket for live game state broadcasting

### 🏗️ Architecture
- **Clean Architecture**: Domain → Engine → Server layers
- **ECS-Inspired Systems**: Movement, Combat, Wave, Projectile, Reward, Lifecycle
- **Data-Driven Design**: YAML configuration for all game balance
- **Repository Pattern**: Pluggable persistence (memory, file, database-ready)
- **Factory Pattern**: Dynamic entity creation from config

### 🎨 Frontend
- **Modern UI**: Glassmorphism design with smooth animations
- **8 Components**: Modular React architecture
- **Tower Selection**: Interactive picker with 3 tower types and stats
- **Game Controls**: Save/Load/Reset with localStorage & file upload/download
- **Responsive**: Works on desktop, tablet, and mobile
- **60 FPS Rendering**: Interpolated canvas animations

### 📊 Observability
- **Prometheus Metrics**: Engine ticks, entities count, performance
- **Structured Logging**: Zap logger with JSON output
- **Request Tracing**: X-Request-ID for debugging
- **pprof Support**: Performance profiling endpoints

---

## 🚀 Quick Start

### Prerequisites
- **Go 1.23+** ([Download](https://go.dev/dl/))
- **Node.js 18+** & npm ([Download](https://nodejs.org/))

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/tower-defense.git
cd tower-defense

# Backend setup
cd backend
go mod download
go build -o server ./cmd/server

# Frontend setup
cd ../frontend
npm install
```


## 📁 Project Structure

### Backend (Clean Architecture)

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── app/                     # Application layer
│   │   ├── router/              # HTTP routing
│   │   │   └── router.go
│   │   └── middleware/          # HTTP middleware
│   │       ├── logging.go
│   │       └── cors.go
│   │
│   ├── infrastructure/          # Infrastructure layer
│   │   ├── websocket/           # WebSocket Hub pattern
│   │   │   └── hub.go
│   │   ├── metrics/             # Prometheus metrics
│   │   │   └── metrics.go
│   │   └── profiling/           # pprof endpoints
│   │       └── pprof.go
│   │
│   ├── game/                    # Game Engine (ECS Architecture)
│   │   ├── ecs/                 # Entity Component System
│   │   │   ├── entity.go        # Entity definitions
│   │   │   ├── factory.go       # Entity factory
│   │   │   └── world.go         # World state
│   │   ├── systems/             # ECS Systems
│   │   │   ├── movement.go      # Movement system
│   │   │   ├── combat.go        # Combat system
│   │   │   ├── projectile.go    # Projectile system
│   │   │   ├── wave.go          # Wave spawning system
│   │   │   ├── reward.go        # Reward system
│   │   │   └── lifecycle.go     # Entity lifecycle
│   │   ├── game.go              # Game instance orchestrator
│   │   ├── manager.go           # Multi-game manager
│   │   └── state.go             # Game state DTOs
│   │
│   ├── repository/              # Persistence layer
│   │   ├── game/                # Game state repository
│   │   │   ├── repository.go
│   │   │   ├── memory.go
│   │   │   └── file.go
│   │   └── config/              # Configuration loader
│   │       ├── loader.go        # YAML config loader
│   │       ├── balance.yaml     # Game balance config
│   │       └── maps.yaml        # Map configurations
│   │
│   ├── config/                  # App configuration
│   │   └── config.go
│   └── logging/                 # Structured logging
│       └── logger.go
│
├── api/
│   └── openapi.yaml             # OpenAPI specification
├── go.mod
└── go.sum
```

### Frontend (Feature-Based Architecture)

```
frontend/
├── src/
│   ├── app/                     # Main application
│   │   ├── App.tsx
│   │   └── App.css
│   │
│   ├── features/                # Feature modules
│   │   ├── game/                # Game feature
│   │   │   └── components/
│   │   │       ├── GameCanvas.tsx
│   │   │       ├── GameOverlay.tsx
│   │   │       └── GameControls.tsx
│   │   ├── hud/                 # HUD feature
│   │   │   └── components/
│   │   │       ├── HUD.tsx
│   │   │       └── TowerSelector.tsx
│   │   ├── maps/                # Maps feature
│   │   │   └── components/
│   │   │       └── MapSelector.tsx
│   │   └── connection/          # Connection status
│   │       └── components/
│   │           └── ConnectionStatus.tsx
│   │
│   ├── shared/                  # Shared resources
│   │   ├── components/
│   │   │   └── ui/              # UI components
│   │   │       ├── Toast.tsx
│   │   │       └── Instructions.tsx
│   │   ├── hooks/               # Shared hooks
│   │   │   └── useToast.ts
│   │   └── config/              # Configuration
│   │       └── config.ts
│   │
│   ├── types.ts                 # TypeScript interfaces
│   └── main.tsx                 # Entry point
│
├── index.html
├── package.json
├── vite.config.ts
└── tsconfig.json
```

---

## 🎮 How to Play

1. **Select a Tower Type** from the side panel (Basic/Sniper/Splash)
2. **Click on the map** to place towers (costs gold)
3. **Defend** against waves of enemies following the blue path
4. **Earn gold** by defeating enemies
5. **Save your progress** anytime with the Save button
6. **Survive as long as possible!** Boss waves every 10 waves

### Tower Types

| Tower | Cost | Damage | Range | Fire Rate | Special |
|-------|------|--------|-------|-----------|---------|
| 🔵 **Basic** | 50 | 10 | 100 | 1.0/s | Balanced |
| 🔴 **Sniper** | 100 | 50 | 200 | 0.5/s | Long-range, high damage |
| 🟠 **Splash** | 75 | 5 | 80 | 2.0/s | Area damage (radius 30) |

### Enemy Types

| Enemy | HP | Speed | Gold | Appears |
|-------|-----|-------|------|---------|
| 🔴 **Basic** | 50 | 1.0x | 10 | Wave 1+ |
| 🔵 **Fast** | 30 | 2.0x | 15 | Wave 6+ (30%) |
| ⚫ **Tank** | 150 | 0.5x | 50 | Wave 11+ (20%) |
| 💜 **Boss** | 500 | 0.75x | 200 | Wave 10, 20, 30... |

---

## 🔧 Configuration

### Game Balance

Edit `backend/internal/repository/config/balance.yaml` to adjust game balance:

```yaml
towers:
  - id: basic
    name: Basic Tower
    cost: 50                # ← Change tower cost
    damage: 10              # ← Change damage
    range: 100              # ← Change range
    fire_rate: 1.0
    projectile_speed: 200

enemies:
  - id: tank
    name: Tank Enemy
    hp: 150                 # ← Change enemy HP
    speed: 0.5              # ← Change speed
    gold_reward: 50         # ← Change gold reward

game:
  starting_gold: 100        # ← Starting resources
  starting_lives: 20
  path:                     # ← Customize enemy path
    - { x: 0, y: 200 }
    - { x: 200, y: 200 }
    # ...
```

**No code changes needed** - just edit YAML and restart! 🎉

---

## 🌐 API Documentation

### REST Endpoints

```
GET  /api/v1/health          # Health check
GET  /api/v1/state           # Current game state
POST /api/v1/tower           # Place tower {x, y, towerType}
POST /api/v1/reset           # Reset game
POST /api/v1/save            # Save game state
POST /api/v1/load            # Load game state

# Multi-room
POST /api/v1/games           # Create new game room
GET  /api/v1/games           # List active rooms

# Legacy endpoints (backward compatibility)
GET  /health                 # Health check
GET  /state                  # Current game state
POST /tower                  # Place tower
POST /reset                  # Reset game
POST /save                   # Save game
POST /load                   # Load game

# Monitoring
GET  /metrics                # Prometheus metrics
GET  /debug/pprof/*          # Performance profiling (if enabled)
```

### WebSocket

```
GET  /ws                     # WebSocket connection
# Receives game state updates ~10 times/second
```

---

## 🛠️ Development

### Adding New Tower Types

1. Edit `backend/internal/repository/config/balance.yaml`:
```yaml
towers:
  - id: laser
    name: Laser Tower
    cost: 200
    damage: 15
    range: 180
    fire_rate: 2.0
    projectile_speed: 800
```

2. Restart backend - **Done!** ✨ (No code changes needed)

### Adding New Enemy Types

1. Edit `backend/internal/repository/config/balance.yaml`:
```yaml
enemies:
  - id: flying
    name: Flying Enemy
    hp: 40
    speed: 1.5
    gold_reward: 20
```

2. Restart backend - **Done!** ✨

### Adding New Maps

Multi-map support can be implemented by:
1. Extending the `balance.yaml` with multiple map configurations
2. Adding map selection UI in frontend
3. Passing selected map to game initialization

---

## 📊 Architecture Highlights

### Backend Architecture (Clean Architecture + ECS)

```
┌─────────────────────────────────────────────────────────┐
│                   Application Layer                      │
│          (Router, Middleware, HTTP Handlers)            │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                Infrastructure Layer                      │
│   WebSocket Hub  │  Metrics  │  Profiling  │  Logging  │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│               Game Manager (Multi-room)                  │
│           (Lifecycle, instance management)              │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│              Game Engine (ECS Pattern)                   │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Systems Layer:                                   │  │
│  │  • Movement    • Combat      • Projectile        │  │
│  │  • Wave        • Reward      • Lifecycle         │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  ECS World (Entity Component System):            │  │
│  │  • Entities (Tower, Enemy, Projectile)           │  │
│  │  • Factory (Dynamic creation from config)        │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                 Repository Layer                         │
│   Config Loader  │  Game State Repository               │
│   (YAML configs) │  (Memory, File, Database-ready)      │
└──────────────────────────────────────────────────────────┘
```

### Key Design Patterns

- **Clean Architecture**: Separation of concerns across layers (App, Infrastructure, Game, Repository)
- **ECS (Entity Component System)**: Game logic split into focused, reusable systems
- **Repository Pattern**: Abstract persistence layer (memory, file, database-ready)
- **Factory Pattern**: Dynamic entity creation from YAML configuration
- **Hub Pattern**: Efficient WebSocket broadcast to multiple clients
- **Feature-Based Frontend**: Modular React architecture with shared components
- **Dependency Injection**: Components receive dependencies via constructors

### Performance Characteristics

- **Backend**: ~60 FPS game loop (16.67ms tick)
- **WebSocket**: ~100ms broadcast interval with adaptive throttling
- **Frontend**: 60 FPS canvas rendering with interpolation
- **Concurrent Games**: Tested with 100+ simultaneous rooms
- **Build Size**: 172 KB (50 KB gzipped)

---

## 📈 Metrics & Monitoring

### Prometheus Metrics

```
# Engine performance
td_engine_ticks_total              # Total game ticks
td_engine_tick_seconds             # Tick duration histogram
td_engine_enemies                  # Current enemy count
td_engine_projectiles              # Current projectile count
td_engine_towers                   # Current tower count

# WebSocket
td_ws_connections                  # Active WebSocket connections

# HTTP
http_requests_total                # Total HTTP requests
http_request_duration_seconds      # Request duration histogram
```

---


### Code Style

- **Go**: Follow [Effective Go](https://go.dev/doc/effective_go)
- **TypeScript/React**: Follow [Airbnb Style Guide](https://github.com/airbnb/javascript/tree/master/react)
---

## 🙏 Acknowledgments

- **Architecture inspiration**: [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- **ECS pattern**: Inspired by Unity ECS and Bevy Engine
- **UI design**: Modern glassmorphism trends
- **Game balance**: Classic tower defense mechanics

---

## 🗺️ Roadmap

### ✅ Completed
- [x] Core gameplay mechanics
- [x] Multiple tower types
- [x] Multiple enemy types
- [x] Save/Load system
- [x] Multi-room support
- [x] Modern frontend UI
- [x] **Clean Architecture refactoring** (Infrastructure, App, Repository layers)
- [x] **Feature-based frontend** organization

### 🚧 In Progress
- [ ] Comprehensive test coverage (>80%)
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] API documentation improvements

### 📅 Planned
- [ ] Docker deployment setup
- [ ] Tower upgrade system
- [ ] Multiple maps (Desert, Jungle, Hell)
- [ ] Power-ups (Freeze, Nuke, Shield)
- [ ] Leaderboard system
- [ ] Achievements
- [ ] Sound effects & music
- [ ] Mobile app (React Native)
- [ ] Backend save/load to file system or database

---

## 🎓 Learning Resources

This project demonstrates:

- **Backend**: Clean Architecture, ECS patterns, WebSocket real-time, Go best practices
- **Frontend**: React hooks, Canvas API, TypeScript, Component architecture
- **DevOps**: Docker, Prometheus, Structured logging
- **Game Dev**: Tower defense mechanics, State management, Entity systems

---

**Made with ❤️ and lots of ☕**

⭐ **Star this repo** if you found it helpful!
