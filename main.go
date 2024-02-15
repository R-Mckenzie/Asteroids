package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/R-Mckenzie/go-engine/engine"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	game := engine.CreateGame(600, 400)

	// SOUNDS
	engine.LoadSound("res/engine.wav", "engine")
	engine.LoadSound("res/explosion.wav", "explosion")
	engine.LoadSound("res/laser.wav", "laser")

	// TEXTURES
	ship := engine.NewTexture("res/playerShip.png")
	bullet := engine.NewTexture("res/bullet.png")
	asteroidBig := engine.NewTexture("res/asteroid_big.png")
	asteroidSmall := engine.NewTexture("res/asteroid_small.png")

	// SPRITES
	playerSprite := engine.NewSprite(96/2, 64/2, 300, 200, 2, ship, nil)

	inGame := &gameScene{
		game:                  game,
		player:                &playerSprite,
		playerDir:             mgl32.Vec2{1, 0},
		bulletTexture:         bullet,
		asteroidBigTexture:    asteroidBig,
		asteroidSmallTexture:  asteroidSmall,
		asteroidTimer:         time.Now(),
		asteroidSpawnInterval: time.Second * 2,
	}

	game.SetScene(inGame)
	menu := &menuScene{game: game, gameScene: inGame}
	game.SetScene(menu)
	game.Run()
}

type menuScene struct {
	game      *engine.Game
	gameScene *gameScene
}

func (m *menuScene) Update() {
	engine.UI.Begin()
	engine.UI.Label("ASTEROIDS", 60, 30, 120, mgl32.Vec4{1, 1, 1, 1})
	if engine.UI.Button(200, 200, 200, 50, "PLAY", mgl32.Vec4{1, 1, 1, 1}) {
		m.game.SetScene(m.gameScene)
		m.gameScene.isPlaying = true
	}

	if engine.UI.Button(200, 300, 200, 50, "QUIT", mgl32.Vec4{1, 1, 1, 1}) {
		m.game.Quit()
	}
	engine.UI.End()
}

type gameScene struct {
	game *engine.Game

	bulletTexture        engine.Texture
	asteroidBigTexture   engine.Texture
	asteroidSmallTexture engine.Texture

	isPlaying bool
	score     int
	playerDir mgl32.Vec2
	playerVel mgl32.Vec2
	player    *engine.Sprite
	asteroids []*entity
	bullets   []*entity

	asteroidTimer         time.Time
	asteroidSpawnInterval time.Duration
}

type entity struct {
	sprite    *engine.Sprite
	direction mgl32.Vec2
}

func (s *gameScene) spawnAsteroid() {
	size := rand.Intn(2)
	side := rand.Intn(4) // 0,1,2,3 starting left going clockwise for which side so spawn in

	var sprite engine.Sprite
	if size == 0 {
		sprite = engine.NewSprite(160/2, 160/2, 0, 0, 0, s.asteroidBigTexture, nil)
	} else {
		sprite = engine.NewSprite(96/2, 96/2, 0, 0, 0, s.asteroidSmallTexture, nil)
	}

	var dir mgl32.Vec2

	switch side {
	case 0: // left
		dir = mgl32.Vec2{1, rand.Float32()}
		sprite.Pos = mgl32.Vec3{-100, float32(rand.Intn(200) + 100), 0}
	case 1: // top
		dir = mgl32.Vec2{rand.Float32(), 1}
		sprite.Pos = mgl32.Vec3{float32(rand.Intn(400) + 100), -100, 0}

	case 2: // right
		dir = mgl32.Vec2{-1, rand.Float32()}
		sprite.Pos = mgl32.Vec3{700, float32(rand.Intn(200) + 100), 0}

	case 3: // bottom
		dir = mgl32.Vec2{rand.Float32(), -1}
		sprite.Pos = mgl32.Vec3{float32(rand.Intn(400) + 100), 500, 0}
	}
	s.asteroids = append(s.asteroids, &entity{sprite: &sprite, direction: dir})
}

var playerSpeed float32 = 0.1
var asteroidSpeed float32 = 2
var rotSpeed float32 = 0.2

func (s *gameScene) Update() {
	if s.isPlaying {
		// MOVEMENT
		if engine.Input.KeyDown(engine.KeyW) {
			s.playerVel = s.playerVel.Add(s.playerDir.Mul(playerSpeed))
			engine.LoopSound("engine", -1)
		}
		if engine.Input.KeyUp(engine.KeyW) {
			engine.StopLoop("engine")
		}

		if engine.Input.KeyDown(engine.KeyS) {
			s.playerVel = s.playerVel.Sub(s.playerDir.Mul(playerSpeed))
			engine.LoopSound("engine", -1)
		}
		if engine.Input.KeyUp(engine.KeyS) {
			engine.StopLoop("engine")
		}

		if engine.Input.KeyDown(engine.KeyA) {
			s.player.Rot[2] -= rotSpeed
			rotMat := mgl32.Rotate2D(-rotSpeed)
			s.playerDir = rotMat.Mul2x1(s.playerDir).Normalize()
		}
		if engine.Input.KeyDown(engine.KeyD) {
			s.player.Rot[2] += rotSpeed
			rotMat := mgl32.Rotate2D(rotSpeed)
			s.playerDir = rotMat.Mul2x1(s.playerDir).Normalize()
		}

		s.player.Pos = s.player.Pos.Add(mgl32.Vec3(s.playerVel.Vec3(0)))

		if s.player.Pos.X() < 0 {
			s.player.Pos[0] = 600
		} else if s.player.Pos.X() > 600 {
			s.player.Pos[0] = 0
		}

		if s.player.Pos.Y() < 0 {
			s.player.Pos[1] = 400
		} else if s.player.Pos.Y() > 400 {
			s.player.Pos[1] = 0
		}

		// SHOOTING
		if engine.Input.KeyOnce(engine.KeySpace) {
			sprite := engine.NewSprite(32, 32, s.player.Pos[0], s.player.Pos[1], 1, s.bulletTexture, nil)
			s.bullets = append(s.bullets, &entity{sprite: &sprite, direction: s.playerDir})
			engine.PlaySound("laser", -1)
		}

		for i, b := range s.bullets {
			if b.sprite.Pos[0] < 0 || b.sprite.Pos[0] > 632 || b.sprite.Pos[1] < 0 || b.sprite.Pos[1] > 432 {
				if i < len(s.bullets)-1 {
					s.bullets = append(s.bullets[:i], s.bullets[i+1:]...)
				} else {
					s.bullets = s.bullets[:len(s.bullets)-1]
				}
				continue
			}

			b.sprite.Pos = b.sprite.Pos.Add(b.direction.Vec3(0).Mul(8))
		}

		if time.Now().Sub(s.asteroidTimer) > s.asteroidSpawnInterval {
			s.spawnAsteroid()
			s.asteroidTimer = time.Now()

			if s.asteroidSpawnInterval > time.Millisecond*100 {
				s.asteroidSpawnInterval -= time.Millisecond * 20
			}
		}

		// ASTEROIDS
		for ia, a := range s.asteroids {
			if collides(*a.sprite, *s.player) {
				s.isPlaying = false
				engine.ClearSounds()
				engine.PlaySound("explosion", 1)
			}

			if a.sprite.Pos[0] < -200 || a.sprite.Pos[0] > 800 || a.sprite.Pos[1] < -200 || a.sprite.Pos[1] > 600 {
				if ia < len(s.asteroids)-1 {
					s.asteroids = append(s.asteroids[:ia], s.asteroids[ia+1:]...)
				} else {
					s.asteroids = s.asteroids[:len(s.asteroids)-1]
				}
				continue
			}

			for ib, b := range s.bullets {
				if collides(*a.sprite, *b.sprite) {
					engine.PlaySound("explosion", -1)
					s.score += 10
					s.bullets = append(s.bullets[:ib], s.bullets[ib+1:]...)
					s.asteroids = append(s.asteroids[:ia], s.asteroids[ia+1:]...)
				}
			}

			a.sprite.Pos = a.sprite.Pos.Add(a.direction.Vec3(0).Mul(asteroidSpeed))
			a.sprite.Rot[2] += 0.1
		}
	}

	engine.Renderer.BeginScene(engine.NewCamera2D(0, 0), mgl32.Vec3{1, 1, 1}, 2)
	engine.Renderer.PushItem(s.player)

	for _, v := range s.bullets {
		engine.Renderer.PushItem(v.sprite)
	}
	for _, v := range s.asteroids {
		engine.Renderer.PushItem(v.sprite)
	}

	engine.UI.Begin()
	engine.UI.Label(fmt.Sprintf("Score: %d", s.score), 10, 10, 32, mgl32.Vec4{1, 1, 1, 1})
	if !s.isPlaying {
		engine.UI.Label(fmt.Sprintf("GAME OVER"), 150, 100, 64, mgl32.Vec4{1, 1, 1, 1})
		if (engine.UI.Button(200, 150, 200, 100, "RESTART", mgl32.Vec4{1, 1, 1, 1})) {
			s.player.Pos = mgl32.Vec3{300, 200, 3}
			s.player.Rot = mgl32.Vec3{}
			s.playerDir = mgl32.Vec2{1, 0}
			s.playerVel = mgl32.Vec2{0, 0}
			s.asteroidSpawnInterval = time.Second * 2
			s.bullets = []*entity{}
			s.asteroids = []*entity{}
			s.score = 0
			s.isPlaying = true
		}
		if (engine.UI.Button(200, 250, 200, 100, "QUIT", mgl32.Vec4{1, 1, 1, 1})) {
			s.game.Quit()
		}
	}
	engine.UI.End()
}

func collides(a, b engine.Sprite) bool {
	ax, ay, aw, ah := a.Pos.X(), a.Pos.Y(), a.Width, a.Height
	bx, by, bw, bh := b.Pos.X(), b.Pos.Y(), b.Width, b.Height

	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}
