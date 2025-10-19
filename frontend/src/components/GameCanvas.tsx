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

      // Path definition
      const path = [
        { x: 0, y: 250 },
        { x: 200, y: 250 },
        { x: 200, y: 100 },
        { x: 400, y: 100 },
        { x: 400, y: 400 },
        { x: 600, y: 400 },
        { x: 600, y: 250 },
        { x: 800, y: 250 },
      ];

      const draw = () => {
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
        path.forEach((p, i) => {
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

        // Get interpolation frames
        const buf = bufferRef.current;
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

        // Draw towers with modern styling
        next.state.towers.forEach((tower) => {
          // Get tower type and color
          const towerType = tower.towerType || 'basic';
          let towerColor = '#4A90E2'; // basic
          let towerSize = 12;
          
          if (towerType === 'sniper') {
            towerColor = '#E24A4A';
            towerSize = 10;
          } else if (towerType === 'splash') {
            towerColor = '#E2A44A';
            towerSize = 14;
          }

          // Range indicator with pulse effect
          const pulseScale = 1 + Math.sin(Date.now() / 500) * 0.05;
          const rangeGradient = ctx.createRadialGradient(
            tower.position.x,
            tower.position.y,
            0,
            tower.position.x,
            tower.position.y,
            tower.range * pulseScale
          );
          const r = parseInt(towerColor.slice(1, 3), 16);
          const g = parseInt(towerColor.slice(3, 5), 16);
          const b = parseInt(towerColor.slice(5, 7), 16);
          rangeGradient.addColorStop(0, `rgba(${r}, ${g}, ${b}, 0.15)`);
          rangeGradient.addColorStop(0.7, `rgba(${r}, ${g}, ${b}, 0.08)`);
          rangeGradient.addColorStop(1, `rgba(${r}, ${g}, ${b}, 0)`);
          ctx.fillStyle = rangeGradient;
          ctx.beginPath();
          ctx.arc(
            tower.position.x,
            tower.position.y,
            tower.range * pulseScale,
            0,
            Math.PI * 2
          );
          ctx.fill();

          // Tower base shadow
          ctx.shadowBlur = 8;
          ctx.shadowColor = 'rgba(0, 0, 0, 0.4)';
          ctx.shadowOffsetY = 2;

          // Tower base
          ctx.fillStyle = towerColor;
          ctx.beginPath();
          ctx.arc(tower.position.x, tower.position.y, towerSize, 0, Math.PI * 2);
          ctx.fill();

          // Tower border
          ctx.shadowBlur = 0;
          ctx.strokeStyle = '#ffffff';
          ctx.lineWidth = 2;
          ctx.beginPath();
          ctx.arc(tower.position.x, tower.position.y, towerSize, 0, Math.PI * 2);
          ctx.stroke();

          // Tower turret with glow
          ctx.shadowBlur = 15;
          ctx.shadowColor = '#2196F3';
          const turretGradient = ctx.createRadialGradient(
            tower.position.x,
            tower.position.y,
            0,
            tower.position.x,
            tower.position.y,
            12
          );
          turretGradient.addColorStop(0, '#64B5F6');
          turretGradient.addColorStop(1, '#1976D2');
          ctx.fillStyle = turretGradient;
          ctx.beginPath();
          ctx.arc(tower.position.x, tower.position.y, 12, 0, Math.PI * 2);
          ctx.fill();

          // Turret highlight
          ctx.shadowBlur = 0;
          ctx.fillStyle = 'rgba(255, 255, 255, 0.4)';
          ctx.beginPath();
          ctx.arc(tower.position.x - 3, tower.position.y - 3, 4, 0, Math.PI * 2);
          ctx.fill();
        });

        ctx.shadowBlur = 0;
        ctx.shadowOffsetY = 0;

        // Draw enemies with modern styling
        enemies.forEach((enemy) => {
          // Get enemy type and properties
          const enemyType = enemy.enemyType || 'basic';
          let enemyColor1 = '#ff6b6b';
          let enemyColor2 = '#c0392b';
          let enemySize = 12;
          
          if (enemyType === 'fast') {
            enemyColor1 = '#4ecdc4';
            enemyColor2 = '#26a69a';
            enemySize = 10;
          } else if (enemyType === 'tank') {
            enemyColor1 = '#95a5a6';
            enemyColor2 = '#7f8c8d';
            enemySize = 16;
          } else if (enemyType === 'boss') {
            enemyColor1 = '#9b59b6';
            enemyColor2 = '#8e44ad';
            enemySize = 20;
          }

          // Enemy shadow
          ctx.shadowBlur = 8;
          ctx.shadowColor = 'rgba(0, 0, 0, 0.4)';
          ctx.shadowOffsetY = 3;

          // Enemy body gradient
          const enemyGradient = ctx.createRadialGradient(
            enemy.position.x - 3,
            enemy.position.y - 3,
            0,
            enemy.position.x,
            enemy.position.y,
            enemySize
          );
          enemyGradient.addColorStop(0, enemyColor1);
          enemyGradient.addColorStop(1, enemyColor2);
          ctx.fillStyle = enemyGradient;
          ctx.beginPath();
          ctx.arc(enemy.position.x, enemy.position.y, enemySize, 0, Math.PI * 2);
          ctx.fill();

          // Enemy border
          ctx.shadowBlur = 0;
          ctx.strokeStyle = 'rgba(0, 0, 0, 0.5)';
          ctx.lineWidth = 2;
          ctx.beginPath();
          ctx.arc(enemy.position.x, enemy.position.y, enemySize, 0, Math.PI * 2);
          ctx.stroke();

          // HP bar background
          const barWidth = enemySize * 2 + 4;
          const barY = enemy.position.y - enemySize - 8;
          ctx.fillStyle = 'rgba(0, 0, 0, 0.6)';
          ctx.fillRect(enemy.position.x - barWidth / 2, barY, barWidth, 5);

          // HP bar fill
          const hpPercent = enemy.hp / enemy.maxHp;
          let hpColor = '#2ecc71'; // green
          if (hpPercent < 0.3) {
            hpColor = '#e74c3c'; // red
          } else if (hpPercent < 0.6) {
            hpColor = '#f1c40f'; // yellow
          }
          
          ctx.fillStyle = hpColor;
          ctx.fillRect(enemy.position.x - barWidth / 2, barY, barWidth * hpPercent, 5);

          // HP bar border
          ctx.strokeStyle = '#fff';
          ctx.lineWidth = 1;
          ctx.strokeRect(enemy.position.x - 18, enemy.position.y - 28, 36, 6);
        });

        ctx.shadowBlur = 0;
        ctx.shadowOffsetY = 0;

        // Draw projectiles with glow trail
        projectiles.forEach((proj) => {
          ctx.shadowBlur = 12;
          ctx.shadowColor = '#f39c12';

          const projGradient = ctx.createRadialGradient(
            proj.position.x,
            proj.position.y,
            0,
            proj.position.x,
            proj.position.y,
            6
          );
          projGradient.addColorStop(0, '#ffd700');
          projGradient.addColorStop(0.5, '#f39c12');
          projGradient.addColorStop(1, '#e67e22');
          ctx.fillStyle = projGradient;
          ctx.beginPath();
          ctx.arc(proj.position.x, proj.position.y, 6, 0, Math.PI * 2);
          ctx.fill();

          // Highlight
          ctx.shadowBlur = 0;
          ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
          ctx.beginPath();
          ctx.arc(proj.position.x - 2, proj.position.y - 2, 2, 0, Math.PI * 2);
          ctx.fill();
        });

        ctx.shadowBlur = 0;

        // Draw game over overlay
        if (hudStateRef.current?.gameOver) {
          ctx.fillStyle = 'rgba(0, 0, 0, 0.85)';
          ctx.fillRect(0, 0, CANVAS_WIDTH, CANVAS_HEIGHT);

          // Game over text with glow
          ctx.shadowBlur = 20;
          ctx.shadowColor = '#e74c3c';
          ctx.fillStyle = '#e74c3c';
          ctx.font = 'bold 56px Arial, sans-serif';
          ctx.textAlign = 'center';
          ctx.fillText('GAME OVER', CANVAS_WIDTH / 2, CANVAS_HEIGHT / 2 - 20);

          ctx.shadowBlur = 10;
          ctx.shadowColor = '#3498db';
          ctx.fillStyle = '#3498db';
          ctx.font = 'bold 28px Arial, sans-serif';
          ctx.fillText(
            `Final Score: ${hudStateRef.current.score}`,
            CANVAS_WIDTH / 2,
            CANVAS_HEIGHT / 2 + 30
          );

          ctx.shadowBlur = 0;
          ctx.fillStyle = '#95a5a6';
          ctx.font = '18px Arial, sans-serif';
          ctx.fillText(
            'Click "Restart" to play again',
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
