import type { GameState } from '../../../types';
import './HUD.css';

interface HUDProps {
  state: GameState | null;
}

export default function HUD({ state }: HUDProps) {
  if (!state) return null;

  return (
    <div className="hud-container">
      <div className="hud-stat wave">
        <div className="hud-icon">🌊</div>
        <div className="hud-content">
          <div className="hud-label">Wave</div>
          <div className="hud-value">{state.wave}</div>
        </div>
      </div>

      <div className="hud-stat gold">
        <div className="hud-icon">💰</div>
        <div className="hud-content">
          <div className="hud-label">Gold</div>
          <div className="hud-value">{state.gold}</div>
        </div>
      </div>

      <div className="hud-stat lives">
        <div className="hud-icon">❤️</div>
        <div className="hud-content">
          <div className="hud-label">Lives</div>
          <div className="hud-value">{state.lives}</div>
        </div>
      </div>

      <div className="hud-stat score">
        <div className="hud-icon">⭐</div>
        <div className="hud-content">
          <div className="hud-label">Score</div>
          <div className="hud-value">{state.score.toLocaleString()}</div>
        </div>
      </div>
      
      <div className="hud-stats-summary">
        <div className="stat-detail">
          <span className="stat-detail-icon">👾</span>
          <span className="stat-detail-label">Enemies:</span>
          <span className="stat-detail-value">{state.enemies.length}</span>
        </div>
        <div className="stat-detail">
          <span className="stat-detail-icon">🗼</span>
          <span className="stat-detail-label">Towers:</span>
          <span className="stat-detail-value">{state.towers.length}</span>
        </div>
      </div>
    </div>
  );
}
