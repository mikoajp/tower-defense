import { TowerType, TowerInfo } from '../types';
import './TowerSelector.css';

interface TowerSelectorProps {
  selectedTower: TowerType;
  onSelectTower: (type: TowerType) => void;
  currentGold: number;
}

export const TOWER_CONFIGS: Record<TowerType, TowerInfo> = {
  basic: {
    type: 'basic',
    name: 'Basic Tower',
    cost: 50,
    damage: 10,
    range: 100,
    fireRate: 1.0,
    description: 'Balanced tower with decent damage and range',
    icon: '🗼',
    color: '#4A90E2'
  },
  sniper: {
    type: 'sniper',
    name: 'Sniper Tower',
    cost: 100,
    damage: 50,
    range: 200,
    fireRate: 0.5,
    description: 'High damage, long range, but slow fire rate',
    icon: '🎯',
    color: '#E24A4A'
  },
  splash: {
    type: 'splash',
    name: 'Splash Tower',
    cost: 75,
    damage: 5,
    range: 80,
    fireRate: 2.0,
    splashRadius: 30,
    description: 'AOE damage - hits multiple enemies at once',
    icon: '💥',
    color: '#E2A44A'
  }
};

export default function TowerSelector({ selectedTower, onSelectTower, currentGold }: TowerSelectorProps) {
  return (
    <div className="tower-selector">
      <h3 className="tower-selector-title">
        <span className="tower-icon">🏗️</span>
        Select Tower
      </h3>
      <div className="tower-grid">
        {(Object.keys(TOWER_CONFIGS) as TowerType[]).map((type) => {
          const tower = TOWER_CONFIGS[type];
          const canAfford = currentGold >= tower.cost;
          const isSelected = selectedTower === type;
          
          return (
            <button
              key={type}
              className={`tower-card ${isSelected ? 'selected' : ''} ${!canAfford ? 'disabled' : ''}`}
              onClick={() => canAfford && onSelectTower(type)}
              disabled={!canAfford}
              style={{ borderColor: isSelected ? tower.color : undefined }}
            >
              <div className="tower-card-header">
                <div className="tower-card-icon">{tower.icon}</div>
                <div className="tower-card-cost">
                  <span className="cost-icon">💰</span>
                  <span className="cost-value">{tower.cost}</span>
                </div>
              </div>
              
              <div className="tower-card-info">
                <div className="tower-card-name">{tower.name}</div>
                
                <div className="tower-card-stats">
                  <div className="stat-item">
                    <span className="stat-icon">⚔️</span>
                    <span className="stat-value">{tower.damage}</span>
                  </div>
                  <div className="stat-item">
                    <span className="stat-icon">📏</span>
                    <span className="stat-value">{tower.range}</span>
                  </div>
                  <div className="stat-item">
                    <span className="stat-icon">⚡</span>
                    <span className="stat-value">{tower.fireRate}/s</span>
                  </div>
                  {tower.splashRadius && (
                    <div className="stat-item">
                      <span className="stat-icon">💥</span>
                      <span className="stat-value">{tower.splashRadius}</span>
                    </div>
                  )}
                </div>
                
                <div className="tower-card-description">{tower.description}</div>
              </div>
              
              {isSelected && <div className="tower-selected-indicator">✓</div>}
              {!canAfford && <div className="tower-disabled-overlay">Not enough gold</div>}
            </button>
          );
        })}
      </div>
      
      <div className="tower-selector-tip">
        <span className="tip-icon">💡</span>
        <span>Click on map to place the selected tower</span>
      </div>
    </div>
  );
}
