package main

import (
	"github.com/R-Mckenzie/go-engine/engine"
	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	game := engine.CreateGame(600, 400)

	menu := menuScene{game: *game}
	game.SetScene(menu)
	game.Run()
	defer game.Quit()
}

type menuScene struct {
	game engine.Game
}

func (m menuScene) Update() {
	engine.Renderer.BeginScene(engine.NewCamera2D(0, 0), mgl32.Vec3{1, 1, 1}, 1)

	engine.UI.Begin()
	engine.UI.Label("ASTEROIDS", 60, 30, 120, mgl32.Vec4{1, 1, 1, 1})
	engine.UI.Button(200, 200, 200, 50, 1, "PLAY", mgl32.Vec4{1, 1, 1, 1})

	if engine.UI.Button(200, 300, 200, 50, 1, "QUIT", mgl32.Vec4{1, 1, 1, 1}) {
		m.game.Quit()
	}
	engine.UI.End()
}
