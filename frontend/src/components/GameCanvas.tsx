import { useEffect, useRef, useImperativeHandle, forwardRef } from 'react';
import type { GameState } from '../types';
import './GameCanvas.css';

const CANVAS_WIDTH = 800;
const CANVAS_HEIGHT = 500;

interface GameCanvasProps {
  onCanvasClick: (x: number, y: number) => void;
  gameOver: boolean;
}

export interface GameCanvasHandle {
  pushState: (state: GameState) => void;
}

type TimedState = { time: number; state: GameState };

const GameCanvas = forwardRef<GameCanvasHandle, GameCanvasProps>(
  ({ onCanvasClick, gameOver }, ref) => {
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const bufferRef = useRef<TimedState[]>([]);
    const rafIdRef = useRef<number | null>(null);
    const hudStateRef = useRef<GameState | null>(null);

    useImperativeHandle(ref, () => ({
      pushState: (state: GameState) => {
        const buf = bufferRef.current;
        buf.push({ time: performance.now(), state });
        if (buf.length > 120) buf.shift();
        hudStateRef.current = state;
      },
    }));

    useEffect(() => {
      if (!canvasRef.current) return;
      const canvas = canvasRef.current;
      const ctx = canvas.getContext('2d');
      if (!ctx) return;

      const INTERPOLATION_DELAY = 120;
      const lerp = (a: number, b: number, t: number) => a + (b - a) * t;

      const draw = () => {
        // Path definition - use path from server state or fallback to default
        const defaultPath = [
          { x: 0, y: 250 },
          { x: 200, y: 250 },
          { x: 200, y: 100 },
          { x: 400, y: 100 },
          { x: 400, y: 400 },
          { x: 600, y: 400 },
          { x: 600, y: 250 },
          { x: 800, y: 250 },
        ];
        
        // Get path from the HUD state (most up-to-date)
        const currentState = hudStateRef.current;
        
        // Debug logging (remove after testing)
        const hasPath = currentState?.path && currentState.path.length > 0;
        if (hasPath && currentState?.path && currentState.path.length !== 8) {
          console.log('Using custom map path with', currentState.path.length, 'points');
        }
        
        const path = hasPath && currentState?.path
          ? currentState.path 
          : defaultPath;
        
        const buf = bufferRef.current;
        // Clear with gradient background
        const bgGradient = ctx.createLinearGradient(0, 0, 0, CANVAS_HEIGHT);
        bgGradient.addColorStop(0, '#1a3a1a');
        bgGradient.addColorStop(1, '#2d5016');
        ctx.fillStyle = bgGradient;
        ctx.fillRect(0, 0, CANVAS_WIDTH, CANVAS_HEIGHT);

        // Draw grid pattern
        ctx.strokeStyle = 'rgba(255, 255, 255, 0.03)';
        ctx.lineWidth = 1;
        for (let x = 0; x < CANVAS_WIDTH; x += 40) {
          ctx.beginPath();
          ctx.moveTo(x, 0);
          ctx.lineTo(x, CANVAS_HEIGHT);
          ctx.stroke();
        }
        for (let y = 0; y < CANVAS_HEIGHT; y += 40) {
          ctx.beginPath();
          ctx.moveTo(0, y);
          ctx.lineTo(CANVAS_WIDTH, y);
          ctx.stroke();
        }

        // Draw path with glow effect
        ctx.shadowBlur = 15;
        ctx.shadowColor = 'rgba(139, 115, 85, 0.5)';
        ctx.strokeStyle = '#8b7355';
        ctx.lineWidth = 44;
        ctx.lineCap = 'round';
        ctx.lineJoin = 'round';
        ctx.beginPath();
        path.forEach((p: { x: number; y: number }, i: number) => {
          if (i === 0) ctx.moveTo(p.x, p.y);
          else ctx.lineTo(p.x, p.y);
        });
        ctx.stroke();

        // Path border
        ctx.shadowBlur = 0;
        ctx.strokeStyle = '#6b5545';
        ctx.lineWidth = 48;
        ctx.stroke();

        ctx.strokeStyle = '#8b7355';
        ctx.lineWidth = 40;
        ctx.stroke();

        // Get interpolation frames (buf already defined above)
        if (buf.length === 0) {
          rafIdRef.current = requestAnimationFrame(draw);
          return;
        }

        const now = performance.now();
        const target = now - INTERPOLATION_DELAY;

        let nextIndex = buf.findIndex((f) => f.time >= target);
        if (nextIndex === -1) nextIndex = buf.length - 1;
        const next = buf[nextIndex];
        const prev = buf[nextIndex - 1] ?? next;

        let alpha = 0;
        if (next.time !== prev.time) {
          alpha = (target - prev.time) / (next.time - prev.time);
          alpha = Math.max(0, Math.min(1, alpha));
        }

        // Interpolate enemies
        const prevEnemiesById = new Map(prev.state.enemies.map((e) => [e.id, e]));
        const enemies = next.state.enemies.map((e) => {
          const p = prevEnemiesById.get(e.id);
          if (!p) return e;
          return {
            ...e,
            position: {
              x: lerp(p.position.x, e.position.x, alpha),
              y: lerp(p.position.y, e.position.y, alpha),
            },
          };
        });

        // Interpolate projectiles
        const prevProjById = new Map(prev.state.projectiles.map((p) => [p.id, p]));
        const projectiles = next.state.projectiles.map((p2) => {
          const p1 = prevProjById.get(p2.id);
          if (!p1) return p2;
          return {
            ...p2,
            position: {
              x: lerp(p1.position.x, p2.position.x, alpha),
              y: lerp(p1.position.y, p2.position.y, alpha),
            },
          };
        });

        // Draw towers with retro pixel styling
        next.state.towers.forEach((tower) => {
          const towerType = tower.towerType || 'basic';
          
          // Range indicator (dotted circle)
          ctx.setLineDash([4, 4]);
          ctx.strokeStyle = 'rgba(255, 215, 0, 0.3)';
          ctx.lineWidth = 1;
          ctx.beginPath();
          ctx.arc(tower.position.x, tower.position.y, tower.range, 0, Math.PI * 2);
          ctx.stroke();
          ctx.setLineDash([]);

          const x = tower.position.x;
          const y = tower.position.y;

          // Draw different tower types
          if (towerType === 'basic') {
            // Basic Tower - Castle turret style
            ctx.fillStyle = '#6b6b6b';
            ctx.fillRect(x - 12, y - 8, 24, 16); // Base
            
            ctx.fillStyle = '#4a4a4a';
            ctx.fillRect(x - 10, y - 12, 6, 8); // Left battlement
            ctx.fillRect(x + 4, y - 12, 6, 8); // Right battlement
            
            ctx.fillStyle = '#8b8b8b';
            ctx.fillRect(x - 12, y - 8, 24, 4); // Highlight
            
            // Cannon
            ctx.fillStyle = '#2c2c2c';
            ctx.fillRect(x - 2, y - 4, 10, 4);
            
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 2;
            ctx.strokeRect(x - 12, y - 8, 24, 16);
            
          } else if (towerType === 'sniper') {
            // Sniper Tower - Tall thin tower with long barrel
            ctx.fillStyle = '#8b4513';
            ctx.fillRect(x - 8, y - 16, 16, 24); // Tall base
            
            ctx.fillStyle = '#654321';
            ctx.fillRect(x - 8, y - 16, 16, 4); // Top
            
            ctx.fillStyle = '#a0522d';
            ctx.fillRect(x - 6, y - 14, 12, 2); // Highlight
            
            // Long sniper barrel
            ctx.fillStyle = '#2c2c2c';
            ctx.fillRect(x - 2, y - 6, 16, 3);
            
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 2;
            ctx.strokeRect(x - 8, y - 16, 16, 24);
            
          } else if (towerType === 'splash') {
            // Splash Tower - Round mortar style
            ctx.fillStyle = '#cd7f32'; // Bronze
            ctx.beginPath();
            ctx.arc(x, y, 14, 0, Math.PI * 2);
            ctx.fill();
            
            ctx.strokeStyle = '#8b6914';
            ctx.lineWidth = 2;
            ctx.beginPath();
            ctx.arc(x, y, 14, 0, Math.PI * 2);
            ctx.stroke();
            
            // Mortar opening
            ctx.fillStyle = '#2c2c2c';
            ctx.beginPath();
            ctx.arc(x, y - 4, 6, 0, Math.PI * 2);
            ctx.fill();
            
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 1;
            ctx.beginPath();
            ctx.arc(x, y - 4, 6, 0, Math.PI * 2);
            ctx.stroke();
            
            // Decorative bands
            ctx.strokeStyle = '#8b6914';
            ctx.lineWidth = 2;
            ctx.beginPath();
            ctx.arc(x, y, 10, 0, Math.PI * 2);
            ctx.stroke();
          }
        });

        ctx.shadowBlur = 0;
        ctx.shadowOffsetY = 0;

        // Draw enemies with retro pixel styling
        enemies.forEach((enemy) => {
          const enemyType = enemy.enemyType || 'basic';
          const x = enemy.position.x;
          const y = enemy.position.y;

          // Draw different enemy types
          if (enemyType === 'basic') {
            // Basic Enemy - Goblin/Orc style
            ctx.fillStyle = '#228b22'; // Green body
            ctx.fillRect(x - 8, y - 8, 16, 16);
            
            // Eyes
            ctx.fillStyle = '#ff0000';
            ctx.fillRect(x - 5, y - 4, 3, 3);
            ctx.fillRect(x + 2, y - 4, 3, 3);
            
            // Teeth
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(x - 4, y + 2, 2, 3);
            ctx.fillRect(x + 2, y + 2, 2, 3);
            
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 2;
            ctx.strokeRect(x - 8, y - 8, 16, 16);
            
          } else if (enemyType === 'fast') {
            // Fast Enemy - Wolf/Beast style
            ctx.fillStyle = '#4169e1'; // Blue body
            
            // Body (horizontal oval shape)
            ctx.beginPath();
            ctx.ellipse(x, y, 12, 8, 0, 0, Math.PI * 2);
            ctx.fill();
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 2;
            ctx.stroke();
            
            // Head
            ctx.fillStyle = '#4169e1';
            ctx.beginPath();
            ctx.arc(x + 8, y, 6, 0, Math.PI * 2);
            ctx.fill();
            ctx.stroke();
            
            // Eyes (glowing)
            ctx.fillStyle = '#ffff00';
            ctx.fillRect(x + 6, y - 2, 2, 2);
            ctx.fillRect(x + 10, y - 2, 2, 2);
            
          } else if (enemyType === 'tank') {
            // Tank Enemy - Armored knight style
            ctx.fillStyle = '#696969'; // Gray armor
            ctx.fillRect(x - 12, y - 12, 24, 24);
            
            ctx.fillStyle = '#2f4f4f'; // Dark accents
            ctx.fillRect(x - 10, y - 10, 20, 4); // Top
            ctx.fillRect(x - 10, y + 6, 20, 4); // Bottom
            
            // Visor
            ctx.fillStyle = '#000000';
            ctx.fillRect(x - 8, y - 2, 16, 4);
            
            // Highlights
            ctx.fillStyle = '#a9a9a9';
            ctx.fillRect(x - 12, y - 12, 4, 24);
            
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 2;
            ctx.strokeRect(x - 12, y - 12, 24, 24);
            
          } else if (enemyType === 'boss') {
            // Boss Enemy - Dragon/Demon style
            ctx.fillStyle = '#8b008b'; // Purple body
            
            // Main body
            ctx.fillRect(x - 16, y - 12, 32, 24);
            
            // Horns
            ctx.fillStyle = '#4b0082';
            ctx.beginPath();
            ctx.moveTo(x - 16, y - 12);
            ctx.lineTo(x - 20, y - 20);
            ctx.lineTo(x - 12, y - 12);
            ctx.fill();
            ctx.beginPath();
            ctx.moveTo(x + 16, y - 12);
            ctx.lineTo(x + 20, y - 20);
            ctx.lineTo(x + 12, y - 12);
            ctx.fill();
            
            // Eyes (glowing red)
            ctx.fillStyle = '#ff0000';
            ctx.fillRect(x - 10, y - 6, 4, 4);
            ctx.fillRect(x + 6, y - 6, 4, 4);
            
            // Mouth
            ctx.fillStyle = '#000';
            ctx.fillRect(x - 8, y + 2, 16, 4);
            
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 3;
            ctx.strokeRect(x - 16, y - 12, 32, 24);
          }

          // HP bar
          const enemySize = enemyType === 'boss' ? 20 : enemyType === 'tank' ? 16 : enemyType === 'fast' ? 12 : 12;
          const barWidth = enemySize * 2 + 4;
          const barY = y - enemySize - 10;
          
          // HP bar background
          ctx.fillStyle = '#000';
          ctx.fillRect(x - barWidth / 2, barY, barWidth, 6);
          
          // HP bar fill
          const hpPercent = enemy.hp / enemy.maxHp;
          let hpColor = '#00ff00'; // Bright green
          if (hpPercent < 0.3) {
            hpColor = '#ff0000'; // Bright red
          } else if (hpPercent < 0.6) {
            hpColor = '#ffff00'; // Bright yellow
          }
          
          ctx.fillStyle = hpColor;
          ctx.fillRect(x - barWidth / 2 + 1, barY + 1, (barWidth - 2) * hpPercent, 4);
          
          // HP bar border
          ctx.strokeStyle = '#fff';
          ctx.lineWidth = 1;
          ctx.strokeRect(x - barWidth / 2, barY, barWidth, 6);
        });

        ctx.shadowBlur = 0;
        ctx.shadowOffsetY = 0;

        // Draw projectiles with retro pixel styling
        projectiles.forEach((proj) => {
          const x = proj.position.x;
          const y = proj.position.y;
          
          // Outer glow
          ctx.fillStyle = 'rgba(255, 215, 0, 0.3)';
          ctx.beginPath();
          ctx.arc(x, y, 8, 0, Math.PI * 2);
          ctx.fill();
          
          // Main projectile body (diamond shape)
          ctx.fillStyle = '#ffd700';
          ctx.beginPath();
          ctx.moveTo(x, y - 5);
          ctx.lineTo(x + 5, y);
          ctx.lineTo(x, y + 5);
          ctx.lineTo(x - 5, y);
          ctx.closePath();
          ctx.fill();
          
          // Inner core
          ctx.fillStyle = '#ffff00';
          ctx.beginPath();
          ctx.moveTo(x, y - 3);
          ctx.lineTo(x + 3, y);
          ctx.lineTo(x, y + 3);
          ctx.lineTo(x - 3, y);
          ctx.closePath();
          ctx.fill();
          
          // Border
          ctx.strokeStyle = '#000';
          ctx.lineWidth = 1;
          ctx.beginPath();
          ctx.moveTo(x, y - 5);
          ctx.lineTo(x + 5, y);
          ctx.lineTo(x, y + 5);
          ctx.lineTo(x - 5, y);
          ctx.closePath();
          ctx.stroke();
        });

        // Draw game over overlay
        if (hudStateRef.current?.gameOver) {
          ctx.fillStyle = 'rgba(0, 0, 0, 0.9)';
          ctx.fillRect(0, 0, CANVAS_WIDTH, CANVAS_HEIGHT);

          // Retro "GAME OVER" text with pixel font effect
          ctx.fillStyle = '#ff0000';
          ctx.font = '48px "Press Start 2P", monospace';
          ctx.textAlign = 'center';
          
          // Text shadow for depth
          ctx.fillStyle = '#000';
          ctx.fillText('GAME OVER', CANVAS_WIDTH / 2 + 4, CANVAS_HEIGHT / 2 - 16);
          
          ctx.fillStyle = '#ff0000';
          ctx.fillText('GAME OVER', CANVAS_WIDTH / 2, CANVAS_HEIGHT / 2 - 20);

          // Score with retro styling
          ctx.fillStyle = '#000';
          ctx.font = '20px "Press Start 2P", monospace';
          ctx.fillText(
            `Score: ${hudStateRef.current.score}`,
            CANVAS_WIDTH / 2 + 3,
            CANVAS_HEIGHT / 2 + 33
          );
          
          ctx.fillStyle = '#ffd700';
          ctx.fillText(
            `Score: ${hudStateRef.current.score}`,
            CANVAS_WIDTH / 2,
            CANVAS_HEIGHT / 2 + 30
          );

          // Restart message
          ctx.fillStyle = '#8b7355';
          ctx.font = '12px "Press Start 2P", monospace';
          ctx.fillText(
            'Press RESTART to try again',
            CANVAS_WIDTH / 2,
            CANVAS_HEIGHT / 2 + 70
          );
        }

        rafIdRef.current = requestAnimationFrame(draw);
      };

      rafIdRef.current = requestAnimationFrame(draw);
      return () => {
        if (rafIdRef.current) cancelAnimationFrame(rafIdRef.current);
      };
    }, []);

    const handleClick = (e: React.MouseEvent<HTMLCanvasElement>) => {
      if (gameOver) return;
      const canvas = canvasRef.current;
      if (!canvas) return;

      const rect = canvas.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      onCanvasClick(x, y);
    };

    return (
      <canvas
        ref={canvasRef}
        width={CANVAS_WIDTH}
        height={CANVAS_HEIGHT}
        onClick={handleClick}
        className="game-canvas"
      />
    );
  }
);

GameCanvas.displayName = 'GameCanvas';

export default GameCanvas;
