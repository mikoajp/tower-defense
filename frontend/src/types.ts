// Shared frontend types matching the backend payloads
export interface Position {
  x: number;
  y: number;
}

export interface Tower {
  id: string;
  towerType: string;
  position: Position;
  range: number;
  damage: number;
  fireRate: number;
  splashRadius?: number;
}

export interface Enemy {
  id: string;
  enemyType: string;
  position: Position;
  hp: number;
  maxHp: number;
  speed: number;
  pathIndex: number;
}

export interface Projectile {
  id: string;
  projectileType: string;
  position: Position;
  target: string;
  speed: number;
  damage: number;
  splashRadius?: number;
}

export interface GameState {
  towers: Tower[];
  enemies: Enemy[];
  projectiles: Projectile[];
  wave: number;
  gold: number;
  lives: number;
  score: number;
  gameOver: boolean;
  path?: Position[];
  mapWidth?: number;
  mapHeight?: number;
}

export type TowerType = 'basic' | 'sniper' | 'splash';

export interface TowerInfo {
  type: TowerType;
  name: string;
  cost: number;
  damage: number;
  range: number;
  fireRate: number;
  splashRadius?: number;
  description: string;
  icon: string;
  color: string;
}
