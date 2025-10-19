import { useEffect, useRef, useState } from 'react';
import { API_URL, WS_URL } from './config';
import type { GameState, TowerType } from './types';
import GameCanvas, { GameCanvasHandle } from './components/GameCanvas';
import HUD from './components/HUD';
import ConnectionStatus from './components/ConnectionStatus';
import GameOverlay from './components/GameOverlay';
import TowerSelector from './components/TowerSelector';
import GameControls from './components/GameControls';
import './App.css';

export default function App() {
  const canvasRef = useRef<GameCanvasHandle>(null);
  const [hudState, setHudState] = useState<GameState | null>(null);
  const [connected, setConnected] = useState(false);
  const [selectedTower, setSelectedTower] = useState<TowerType>('basic');

  useEffect(() => {
    let ws: WebSocket | null = null;
    let reconnectTimeout: number | null = null;
    let isComponentMounted = true;

    const connect = () => {
      if (!isComponentMounted) return;

      if (ws && ws.readyState !== WebSocket.CLOSED) {
        ws.close();
      }

      ws = new WebSocket(WS_URL);

      ws.onopen = () => {
        if (!isComponentMounted) {
          ws?.close();
          return;
        }
        console.log('🎮 Connected to game server');
        setConnected(true);
      };

      ws.onmessage = (event) => {
        if (!isComponentMounted) return;

        try {
          const raw = JSON.parse(event.data);
          const state: GameState = {
            ...raw,
            towers: raw.towers ?? [],
            enemies: raw.enemies ?? [],
            projectiles: raw.projectiles ?? [],
          };

          canvasRef.current?.pushState(state);
          setHudState(state);
        } catch (error) {
          console.error('Error parsing game state:', error);
        }
      };

      ws.onclose = () => {
        if (!isComponentMounted) return;

        console.log('❌ Disconnected from server');
        setConnected(false);

        if (isComponentMounted) {
          reconnectTimeout = window.setTimeout(() => {
            if (isComponentMounted) {
              console.log('🔄 Attempting to reconnect...');
              connect();
            }
          }, 2000);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    };

    connect();

    return () => {
      isComponentMounted = false;

      if (reconnectTimeout) {
        window.clearTimeout(reconnectTimeout);
      }

      if (ws) {
        ws.onclose = null;
        ws.close();
      }
    };
  }, []);

  const handleCanvasClick = async (x: number, y: number) => {
    if (!hudState || hudState.gameOver) return;

    try {
      const response = await fetch(`${API_URL}/tower`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ x, y, towerType: selectedTower }),
      });

      if (!response.ok) {
        const data = await response.json();
        alert(data.error || 'Failed to place tower');
      }
    } catch (error) {
      console.error('Error placing tower:', error);
    }
  };

  const handleRestart = async () => {
    try {
      await fetch(`${API_URL}/reset`, { method: 'POST' });
    } catch (e) {
      console.error('Failed to reset game', e);
    }
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <h1 className="app-title">
          <span className="title-icon">🏰</span>
          <span>Tower Defense</span>
        </h1>
        <ConnectionStatus connected={connected} />
      </header>

      <main className="app-main">
        <div className="game-layout">
          <div className="left-panel">
            {hudState && <HUD state={hudState} />}
            
            <TowerSelector
              selectedTower={selectedTower}
              onSelectTower={setSelectedTower}
              currentGold={hudState?.gold ?? 0}
            />
            
            <GameControls />
          </div>

          <div className="game-area">
            <GameOverlay show={hudState?.gameOver ?? false} onRestart={handleRestart} />
            
            <div className="canvas-wrapper">
              <GameCanvas
                ref={canvasRef}
                onCanvasClick={handleCanvasClick}
                gameOver={hudState?.gameOver ?? false}
              />
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
