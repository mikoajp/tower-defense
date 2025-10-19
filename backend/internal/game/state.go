package game

import "errors"

var (
	ErrNotEnoughGold    = errors.New("not enough gold")
	ErrInvalidPlacement = errors.New("invalid tower placement")
	ErrGameNotFound     = errors.New("game not found")
)

// GameStateSnapshot represents a snapshot of the game state for serialization
type GameStateSnapshot struct {
	Towers      []TowerDTO      `json:"towers"`
	Enemies     []EnemyDTO      `json:"enemies"`
	Projectiles []ProjectileDTO `json:"projectiles"`
	Wave        int             `json:"wave"`
	Gold        int             `json:"gold"`
	Lives       int             `json:"lives"`
	Score       int             `json:"score"`
	GameOver    bool            `json:"gameOver"`
	Path        []PosDTO        `json:"path"`
	MapWidth    int             `json:"mapWidth"`
	MapHeight   int             `json:"mapHeight"`
}

// TowerDTO is the data transfer object for towers
type TowerDTO struct {
	ID           string  `json:"id"`
	Type         string  `json:"towerType"`
	Position     PosDTO  `json:"position"`
	Range        float64 `json:"range"`
	Damage       int     `json:"damage"`
	FireRate     float64 `json:"fireRate"`
	SplashRadius float64 `json:"splashRadius,omitempty"`
}

// EnemyDTO is the data transfer object for enemies
type EnemyDTO struct {
	ID        string  `json:"id"`
	Type      string  `json:"enemyType"`
	Position  PosDTO  `json:"position"`
	HP        int     `json:"hp"`
	MaxHP     int     `json:"maxHp"`
	Speed     float64 `json:"speed"`
	PathIndex int     `json:"pathIndex"`
}

// ProjectileDTO is the data transfer object for projectiles
type ProjectileDTO struct {
	ID           string  `json:"id"`
	Type         string  `json:"projectileType"`
	Position     PosDTO  `json:"position"`
	Target       string  `json:"target"`
	Speed        float64 `json:"speed"`
	Damage       int     `json:"damage"`
	SplashRadius float64 `json:"splashRadius,omitempty"`
}

// PosDTO is the data transfer object for positions
type PosDTO struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Convert ECS entities to DTOs
func (g *Game) convertTowers() []TowerDTO {
	towers := g.world.GetTowers()
	dtos := make([]TowerDTO, 0, len(towers))
	
	for _, t := range towers {
		dtos = append(dtos, TowerDTO{
			ID:           t.ID,
			Type:         t.TowerType,
			Position:     PosDTO{X: t.Position.X, Y: t.Position.Y},
			Range:        t.Range,
			Damage:       t.Damage,
			FireRate:     t.FireRate,
			SplashRadius: t.SplashRadius,
		})
	}
	
	return dtos
}

func (g *Game) convertEnemies() []EnemyDTO {
	enemies := g.world.GetEnemies()
	dtos := make([]EnemyDTO, 0, len(enemies))
	
	for _, e := range enemies {
		dtos = append(dtos, EnemyDTO{
			ID:        e.ID,
			Type:      e.EnemyType,
			Position:  PosDTO{X: e.Position.X, Y: e.Position.Y},
			HP:        e.HP,
			MaxHP:     e.MaxHP,
			Speed:     e.Speed,
			PathIndex: e.PathIndex,
		})
	}
	
	return dtos
}

func (g *Game) convertProjectiles() []ProjectileDTO {
	projectiles := g.world.GetProjectiles()
	dtos := make([]ProjectileDTO, 0, len(projectiles))
	
	for _, p := range projectiles {
		dtos = append(dtos, ProjectileDTO{
			ID:           p.ID,
			Type:         p.ProjectileType,
			Position:     PosDTO{X: p.Position.X, Y: p.Position.Y},
			Target:       p.Target,
			Speed:        p.Speed,
			Damage:       p.Damage,
			SplashRadius: p.SplashRadius,
		})
	}
	
	return dtos
}
