package game

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Item struct {
	gravity  float32
	itemID   int
	position rl.Vector2
	rotation float32
	collided bool
}

type Boundaries struct {
	bottom float32
	top    float32
	width  float32
}

type Spawner struct {
	boundaries Boundaries
	clock      float32
	gravity    float32
	idPoolSize int
	items      []Item
	spawnLimit int
	timeOut    float32
}

func NewSpawner(duration float32, gravity float32, idPoolSize int, spawnLimit int, spawnBoundaries Boundaries) Spawner {
	return Spawner{
		boundaries: spawnBoundaries,
		clock:      duration,
		gravity:    gravity,
		idPoolSize: idPoolSize,
		items:      []Item{},
		spawnLimit: spawnLimit,
		timeOut:    0.0,
	}
}

func UpdateSpawner(spawner *Spawner, mover func(*Item, float32) bool, delta float32) {
	if !canSpawn(spawner) {
		spawner.clock = 0.0
	} else {
		spawnItem(spawner, delta)
	}

	if len(spawner.items) == 0 {
		return
	}

	shouldContinue := true

	for i := range spawner.items {
		shouldContinue = mover(&spawner.items[i], delta)

		if !shouldContinue {
			break
		}
	}

	if shouldContinue {
		reSpawnItems(spawner)
	}
}

func canSpawn(spawner *Spawner) bool {
	return len(spawner.items) < spawner.spawnLimit
}

func spawnItem(spawner *Spawner, dt float32) {
	spawner.clock += dt

	if spawner.clock >= spawner.timeOut {
		spawner.clock = 0.0

		spawner.items = append(spawner.items, makeItem(spawner.boundaries.width, spawner.boundaries.top, spawner.idPoolSize, spawner.gravity))
	}
}

func makeItem(boundaryX float32, boundaryY float32, idPoolSize int, gravityRange float32) Item {
	spawnX := float32(rand.Intn(int(boundaryX))) + float32(rand.Intn(100))
	spawnY := boundaryY - float32(rand.Intn(100))
	gravity := float32(RandRange(100, int(gravityRange)))

	return Item{
		position: rl.NewVector2(spawnX, spawnY),
		itemID:   rand.Intn(idPoolSize),
		gravity:  gravity,
	}
}

func reSpawnItems(spawner *Spawner) {
	for i, item := range spawner.items {
		if item.position.Y >= spawner.boundaries.bottom || item.collided {
			spawner.items[i] = makeItem(spawner.boundaries.width, spawner.boundaries.top, spawner.idPoolSize, spawner.gravity)
		}

	}
}
