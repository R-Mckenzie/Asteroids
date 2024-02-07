package main

import (
	_ "flag"
	"fmt"
	_ "log"
	_ "os"
	_ "runtime"
	_ "runtime/pprof"

	"github.com/R-Mckenzie/go-engine/engine"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	game := engine.CreateGame(600, 400)

	ship := engine.NewTexture("res/playerShip.png")
	bullet := engine.NewTexture("res/bullet.png")
	asteroidBig := engine.NewTexture("res/asteroid_big.png")
	asteroidSmall := engine.NewTexture("res/asteroid_small.png")

	playerSprite := engine.NewSprite(96/2, 64/2, 300, 200, 2, ship, nil)

	inGame := &gameScene{
		game:                 game,
		player:               &playerSprite,
		bulletTexture:        bullet,
		asteroidBigTexture:   asteroidBig,
		asteroidSmallTexture: asteroidSmall,
	}

	game.SetScene(inGame)
	menu := &menuScene{game: game, gameScene: inGame}
	game.SetScene(menu)
	game.Run()
	defer game.Quit()
}

type menuScene struct {
	game      *engine.Game
	gameScene *gameScene
}

func (m *menuScene) Update() {
	if engine.Input.KeyUp(engine.KeyF) {
		fmt.Println("FFF")
	}
	engine.Renderer.BeginScene(engine.NewCamera2D(0, 0), mgl32.Vec3{1, 1, 1}, 1)

	engine.UI.Begin()
	engine.UI.Label("ASTEROIDS", 60, 30, 120, mgl32.Vec4{1, 1, 1, 1})
	if engine.UI.Button(200, 200, 200, 50, 1, "PLAY", mgl32.Vec4{1, 1, 1, 1}) {
		m.game.SetScene(m.gameScene)
	}

	if engine.UI.Button(200, 300, 200, 50, 1, "QUIT", mgl32.Vec4{1, 1, 1, 1}) {
		// m.game.Quit()
	}
	engine.UI.End()
}

type gameScene struct {
	game *engine.Game

	bulletTexture        engine.Texture
	asteroidBigTexture   engine.Texture
	asteroidSmallTexture engine.Texture

	score     int
	player    *engine.Sprite
	asteroids []engine.Sprite
	bullets   []*bullet
}

var speed float32 = 10

type bullet struct {
	sprite    *engine.Sprite
	direction mgl32.Vec2
}

func (s *gameScene) Update() {

	// MOVEMENT
	if engine.Input.KeyDown(engine.KeyW) {
		s.player.Pos[1] -= speed
	}
	if engine.Input.KeyDown(engine.KeyS) {
		s.player.Pos[1] += speed
	}
	if engine.Input.KeyDown(engine.KeyA) {
		s.player.Pos[0] -= speed
	}
	if engine.Input.KeyDown(engine.KeyD) {
		s.player.Pos[0] += speed
	}

	// SHOOTING
	if engine.Input.KeyOnce(engine.KeySpace) {
		sprite := engine.NewSprite(32, 32, s.player.Pos[0], s.player.Pos[1], 1, s.bulletTexture, nil)
		s.bullets = append(s.bullets,
			&bullet{sprite: &sprite, direction: mgl32.Vec2{1, 0}})
	}

	for i, b := range s.bullets {
		if b.sprite.Pos[0] < 0 || b.sprite.Pos[0] > 632 || b.sprite.Pos[1] < 0 || b.sprite.Pos[1] > 432 {
			s.bullets = append(s.bullets[:i], s.bullets[i+1:]...)
			continue
		}
		b.sprite.Pos = b.sprite.Pos.Add(b.direction.Vec3(0).Mul(5))
	}

	engine.Renderer.BeginScene(engine.NewCamera2D(0, 0), mgl32.Vec3{1, 1, 1}, 2)
	engine.Renderer.PushItem(s.player)

	for _, v := range s.bullets {
		engine.Renderer.PushItem(v.sprite)
	}
	for _, v := range s.asteroids {
		engine.Renderer.PushItem(v)
	}

	engine.UI.Begin()
	engine.UI.Label(fmt.Sprintf("Score: %d", s.score), 10, 10, 32, mgl32.Vec4{1, 1, 1, 1})
	engine.UI.End()
}
