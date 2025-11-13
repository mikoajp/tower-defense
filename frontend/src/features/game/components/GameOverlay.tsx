import './GameOverlay.css';

interface GameOverlayProps {
  show: boolean;
  onRestart: () => void;
}

export default function GameOverlay({ show, onRestart }: GameOverlayProps) {
  if (!show) return null;

  return (
    <div className="game-overlay-container">
      <button onClick={onRestart} className="restart-button">
        <span className="restart-icon">🔄</span>
        <span>Restart Game</span>
      </button>
    </div>
  );
}
