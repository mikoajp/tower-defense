# 🎮 Tower Defense Frontend - Modern UI

Modern, component-based React + TypeScript frontend with advanced styling and animations.

## ✨ Features

### 🎨 Modern Design
- **Glassmorphism UI** - Frosted glass effects with backdrop blur
- **Smooth Animations** - Fade-ins, slides, pulses, and hover effects
- **Gradient Backgrounds** - Dynamic color schemes
- **Responsive Design** - Mobile-friendly layout

### 🧩 Component Architecture
```
src/
├── App.tsx                     # Main application container
├── App.css                     # Global styles
├── components/                 # 8 React components
│   ├── GameCanvas.tsx          # Canvas rendering with interpolation
│   ├── GameCanvas.css          # Canvas styling
│   ├── HUD.tsx                 # Game statistics display
│   ├── HUD.css                 # HUD grid layout
│   ├── TowerSelector.tsx       # Tower type picker with stats
│   ├── TowerSelector.css       # Tower selector styling
│   ├── GameControls.tsx        # Save/Load/Reset controls
│   ├── GameControls.css        # Game controls styling
│   ├── ConnectionStatus.tsx    # WebSocket connection indicator
│   ├── ConnectionStatus.css    # Connection status styling
│   ├── GameOverlay.tsx         # Game over restart button
│   ├── GameOverlay.css         # Overlay animations
│   ├── Instructions.tsx        # How to play guide
│   └── Instructions.css        # Instructions layout
├── types.ts                    # TypeScript interfaces
└── config.ts                   # API/WS configuration
```

## 🎯 Component Details

### **App.tsx**
Main application container that:
- Manages WebSocket connection lifecycle
- Coordinates state between components
- Handles tower placement API calls
- Provides game restart functionality

### **GameCanvas.tsx**
Advanced canvas rendering with:
- **60 FPS interpolation** for smooth animations
- **Gradient path rendering** with glow effects
- **Pulsing tower ranges** with radial gradients
- **HP bars** with color-coded health states
- **Projectile trails** with glow effects
- **Game over overlay** rendered on canvas

### **HUD.tsx**
Modern statistics display featuring:
- **Grid layout** - 4-column responsive grid
- **Glassmorphism cards** - Translucent backgrounds
- **Color-coded stats** - Each stat has unique gradient
- **Hover animations** - Lift effect on hover

### **ConnectionStatus.tsx**
Connection indicator with:
- **Animated spinner** - CSS rotation animation
- **Pulse effect** - Box shadow animation
- **Auto-hide** when connected

### **GameOverlay.tsx**
Game restart control with:
- **Slide-down animation** on appearance
- **Rotating restart icon** - Continuous rotation
- **Gradient button** with hover effects

### **Instructions.tsx**
Player guide featuring:
- **2-column grid layout** (responsive to 1 column on mobile)
- **Icon-based instructions** with hover effects
- **Tip section** with special styling

## 🎨 Design System

### **Color Palette**
```css
/* Background */
Primary: linear-gradient(135deg, #0f2027, #203a43, #2c5364)

/* UI Elements */
Glass: rgba(255, 255, 255, 0.08) - 0.15
Borders: rgba(255, 255, 255, 0.1) - 0.2

/* Accent Colors */
Wave: #3498db (Blue)
Gold: #f1c40f (Yellow)
Lives: #e74c3c (Red)
Score: #2ecc71 (Green)
Tower: #2196F3 (Light Blue)
Enemy: #e74c3c → #c0392b (Red gradient)
Projectile: #ffd700 → #e67e22 (Gold gradient)
```

### **Typography**
```css
Headings: -apple-system, BlinkMacSystemFont, 'Segoe UI'
Font Weights: 500 (medium), 600 (semibold), 700 (bold), 800 (extrabold)
```

### **Effects**
- **Glassmorphism**: `backdrop-filter: blur(10px)`
- **Shadows**: Multi-layer box-shadows for depth
- **Transitions**: 0.3s ease for smooth interactions
- **Animations**: CSS keyframes for complex effects

## 📱 Responsive Design

### **Breakpoints**
```css
Mobile: max-width: 768px
- HUD: 2x2 grid instead of 1x4
- Instructions: Single column
- Font sizes reduced
- Padding adjusted
```

## 🚀 Running the Application

### **Development**
```bash
npm install
npm run dev
# Opens http://localhost:3000 (configured in vite.config.ts)
```

### **Production Build**
```bash
npm run build
npm run preview
```

### **Environment Variables**
```bash
# .env.local
VITE_API_URL=http://localhost:8080
```

## 🎮 User Interactions

### **Tower Placement**
1. Click canvas → Sends POST `/tower` with {x, y}
2. Visual feedback on hover (crosshair cursor)
3. Error alerts for insufficient gold

### **Game Restart**
1. Click restart button → Sends POST `/reset`
2. Button animates with rotating icon
3. Canvas refreshes automatically

### **WebSocket Updates**
1. Receives game state ~10 times/second
2. Interpolates entity positions for smooth 60fps
3. Auto-reconnects after 2 seconds on disconnect

## 🏗️ Architecture Patterns

### **Separation of Concerns**
- **App.tsx**: State management & API communication
- **GameCanvas.tsx**: Pure rendering logic
- **HUD/Instructions**: Presentational components
- **CSS Modules**: Scoped styles per component

### **Performance Optimizations**
- **useRef** for canvas and RAF to avoid re-renders
- **Interpolation buffer** (120 frames) for smooth animations
- **RequestAnimationFrame** for optimal rendering
- **Memoization** via React.memo (can be added)

### **Type Safety**
- Full TypeScript coverage
- Interfaces matching backend domain types
- Type-safe API calls and state management

## 🎨 Animation Showcase

### **Entry Animations**
```css
Header: fadeInDown 0.6s
HUD: fadeIn 0.8s (delay 0.2s)
Canvas: fadeIn 1s (delay 0.4s)
Instructions: fadeIn 1s (delay 0.6s)
```

### **Interaction Animations**
```css
HUD Cards: translateY on hover
Tower Ranges: Pulsing scale (sin wave)
Restart Icon: Continuous rotation
Connection Status: Box shadow pulse
```

## 📊 Component Props Interface

```typescript
// GameCanvas
interface GameCanvasProps {
  onCanvasClick: (x: number, y: number) => void;
  gameOver: boolean;
}

// HUD
interface HUDProps {
  state: GameState | null;
}

// ConnectionStatus
interface ConnectionStatusProps {
  connected: boolean;
}

// GameOverlay
interface GameOverlayProps {
  show: boolean;
  onRestart: () => void;
}
```

## 🔧 Future Enhancements

- [x] Tower selection menu with different types (COMPLETED)
- [x] Game controls with save/load (COMPLETED)
- [ ] Upgrade system UI
- [ ] Wave preview/countdown
- [ ] Leaderboard component
- [ ] Sound effects toggle
- [ ] Fullscreen mode
- [ ] Touch controls optimization
- [ ] Particle effects library integration

## 📝 CSS Best Practices

1. **BEM-like naming** for clarity
2. **CSS custom properties** for theming (can be added)
3. **Mobile-first** responsive design
4. **Accessibility** - Focus states and ARIA labels (can be improved)
5. **Performance** - GPU-accelerated transforms

---

**Built with ❤️ using React, TypeScript, and modern CSS**
