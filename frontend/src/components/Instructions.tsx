import './Instructions.css';

export default function Instructions() {
  return (
    <div className="instructions-container">
      <h3 className="instructions-title">📖 How to Play</h3>
      <div className="instructions-grid">
        <div className="instruction-item">
          <div className="instruction-icon">🎯</div>
          <div className="instruction-text">
            <strong>Place Towers:</strong> Click anywhere on the canvas to place a tower (costs 50 gold)
          </div>
        </div>
        <div className="instruction-item">
          <div className="instruction-icon">🔫</div>
          <div className="instruction-text">
            <strong>Auto-Fire:</strong> Towers automatically shoot enemies within range
          </div>
        </div>
        <div className="instruction-item">
          <div className="instruction-icon">👾</div>
          <div className="instruction-text">
            <strong>Defend:</strong> Survive the waves - don't let enemies reach the end!
          </div>
        </div>
        <div className="instruction-item">
          <div className="instruction-icon">💰</div>
          <div className="instruction-text">
            <strong>Earn Gold:</strong> Defeating enemies rewards you with gold for more towers
          </div>
        </div>
      </div>
      <div className="instructions-tip">
        <span className="tip-icon">💡</span>
        <span>Tip: Place towers strategically at corners where enemies slow down!</span>
      </div>
    </div>
  );
}
