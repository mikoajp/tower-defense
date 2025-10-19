import { useEffect, useRef, useState } from 'react';
import { API_URL, WS_URL } from './config';
import type { GameState, TowerType } from './types';
import GameCanvas, { GameCanvasHandle } from './components/GameCanvas';
import HUD from './components/HUD';
import ConnectionStatus from './components/ConnectionStatus';
import GameOverlay from './components/GameOverlay';
import TowerSelector from './components/TowerSelector';
import GameControls from './components/GameControls';
import MapSelector from './components/MapSelector';
import Toast from './components/Toast';
import { useToast } from './hooks/useToast';
import './App.css';

export default function App() {
  const canvasRef = useRef<GameCanvasHandle>(null);
  const [hudState, setHudState] = useState<GameState | null>(null);
  const [connected, setConnected] = useState(false);
  const [selectedTower, setSelectedTower] = useState<TowerType>('basic');
  const { toasts, removeToast, showError, showSuccess, showWarning, showInfo } = useToast();

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
        showSuccess('Connected to game server!');
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
          showError('Failed to parse game data');
        }
      };

      ws.onclose = () => {
        if (!isComponentMounted) return;

        console.log('❌ Disconnected from server');
        setConnected(false);
        showWarning('Disconnected from server. Reconnecting...');

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
        showError('Connection error occurred');
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
        const errorMsg = data.error || 'Failed to place tower';
        
        if (errorMsg.includes('not enough gold')) {
          showError('Not enough gold!');
        } else if (errorMsg.includes('invalid')) {
          showError('Cannot place tower here!');
        } else {
          showError(errorMsg);
        }
      } else {
        showSuccess(`${selectedTower.charAt(0).toUpperCase() + selectedTower.slice(1)} tower placed!`);
      }
    } catch (error) {
      console.error('Error placing tower:', error);
      showError('Network error: Failed to place tower');
    }
  };

  const handleRestart = async () => {
    try {
      const response = await fetch(`${API_URL}/reset`, { method: 'POST' });
      if (response.ok) {
        showInfo('Game restarted!');
      } else {
        showError('Failed to restart game');
      }
    } catch (e) {
      console.error('Failed to reset game', e);
      showError('Network error: Failed to restart game');
    }
  };

  const handleMapChange = (_mapId: string, mapName: string) => {
    showSuccess(`Map changed to: ${mapName}`);
  };

  return (
    <div className="app-container">
      {/* Toast Notifications */}
      <div className="toast-container">
        {toasts.map((toast) => (
          <Toast
            key={toast.id}
            message={toast.message}
            type={toast.type}
            onClose={() => removeToast(toast.id)}
          />
        ))}
      </div>

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
            
            <MapSelector onMapChange={handleMapChange} />
            
            <TowerSelector
              selectedTower={selectedTower}
              onSelectTower={setSelectedTower}
              currentGold={hudState?.gold ?? 0}
            />
            
            <GameControls 
              showSuccess={showSuccess}
              showError={showError}
            />
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
