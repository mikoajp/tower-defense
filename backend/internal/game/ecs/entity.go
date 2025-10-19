package ecs

import "time"

// EntityType represents the type of game entity
type EntityType string

const (
	EntityTypeTower      EntityType = "tower"
	EntityTypeEnemy      EntityType = "enemy"
	EntityTypeProjectile EntityType = "projectile"
)

// Entity is the base interface for all game entities
type Entity interface {
	GetID() string
	GetType() EntityType
	GetPosition() Position
	SetPosition(Position)
	Update(dt float64)
	IsAlive() bool
}

// Position represents a 2D point
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// BaseEntity provides common entity functionality
type BaseEntity struct {
	ID       string
	Type     EntityType
	Position Position
	Alive    bool
}

func (e *BaseEntity) GetID() string {
	return e.ID
}

func (e *BaseEntity) GetType() EntityType {
	return e.Type
}

func (e *BaseEntity) GetPosition() Position {
	return e.Position
}

func (e *BaseEntity) SetPosition(pos Position) {
	e.Position = pos
}

func (e *BaseEntity) IsAlive() bool {
	return e.Alive
}

// TowerEntity represents a defense tower
type TowerEntity struct {
	BaseEntity
	TowerType    string    `json:"towerType"`
	Range        float64   `json:"range"`
	Damage       int       `json:"damage"`
	FireRate     float64   `json:"fireRate"`
	SplashRadius float64   `json:"splashRadius,omitempty"`
	LastShot     time.Time `json:"-"`
}

func (t *TowerEntity) Update(dt float64) {
	// Towers are stationary, no update needed
}

func (t *TowerEntity) CanShoot() bool {
	elapsed := time.Since(t.LastShot).Seconds()
	return elapsed >= 1.0/t.FireRate
}

func (t *TowerEntity) Shoot() {
	t.LastShot = time.Now()
}

// EnemyEntity represents an enemy
type EnemyEntity struct {
	BaseEntity
	EnemyType string  `json:"enemyType"`
	HP        int     `json:"hp"`
	MaxHP     int     `json:"maxHp"`
	Speed     float64 `json:"speed"`
	PathIndex int     `json:"pathIndex"`
	GoldReward  int   `json:"-"`
	ScoreReward int   `json:"-"`
}

func (e *EnemyEntity) Update(dt float64) {
	// Movement handled by MovementSystem
}

func (e *EnemyEntity) TakeDamage(damage int) {
	e.HP -= damage
	if e.HP <= 0 {
		e.HP = 0
		e.Alive = false
	}
}

func (e *EnemyEntity) GetHealthPercent() float64 {
	if e.MaxHP == 0 {
		return 0
	}
	return float64(e.HP) / float64(e.MaxHP)
}

// ProjectileEntity represents a projectile
type ProjectileEntity struct {
	BaseEntity
	ProjectileType string  `json:"projectileType"`
	Target         string  `json:"target"`
	Speed          float64 `json:"speed"`
	Damage         int     `json:"damage"`
	SplashRadius   float64 `json:"splashRadius,omitempty"`
}

func (p *ProjectileEntity) Update(dt float64) {
	// Movement handled by ProjectileSystem
}

// Damageable represents entities that can take damage
type Damageable interface {
	TakeDamage(damage int)
	GetHealthPercent() float64
}

// Shooter represents entities that can shoot
type Shooter interface {
	CanShoot() bool
	Shoot()
}
