package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	modelCube = iota
	modelPenger
	modelLength
)

const (
	viewModeVertices = iota
	viewModeFaces
	viewModeBoth
	viewModeLength
)

const (
	width         = 800
	height        = 800
	vertexSize    = 7
	lineThickness = 1
	fps           = 60
)

var (
	backgroundColor         = rl.Black
	verticesColor           = rl.Red
	facesColor              = rl.Green
	delta           float32 = 1
	angle           float32
	currentModel    = modelCube
	currentViewMode = viewModeBoth
	vertices        = CubeVertices
	faces           = CubeFaces
)

func main() {
	rl.InitWindow(width, height, "Tsoding 3D")
	rl.SetTargetFPS(fps)

	for !rl.WindowShouldClose() {
		handleKeyboard()

		rl.BeginDrawing()

		angle += math.Pi * rl.GetFrameTime()
		rl.ClearBackground(backgroundColor)

		if currentViewMode == viewModeVertices || currentViewMode == viewModeBoth {
			for _, vertex := range vertices {
				vertex.RotateXY(angle).TranslateZ(delta).Project().Screen().Point()
			}
		}

		if currentViewMode == viewModeFaces || currentViewMode == viewModeBoth {
			for _, face := range faces {
				for index := 0; index < len(face); index++ {
					a := vertices[face[index]]
					b := vertices[face[(index+1)%len(face)]]

					v1 := a.RotateXY(angle).TranslateZ(delta).Project().Screen()
					v2 := b.RotateXY(angle).TranslateZ(delta).Project().Screen()

					v1.Line(v2)
				}
			}
		}

		rl.DrawText("V - [V]iew mode - Switch between vertices/faces/both", 10, 10, 20, rl.Blue)
		rl.DrawText("M - [M]odel - Switch between cube/penger model", 10, 30, 20, rl.Blue)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func handleKeyboard() {
	if rl.IsKeyPressed(rl.KeyV) {
		currentViewMode = (currentViewMode + 1) % viewModeLength
	}
	if rl.IsKeyPressed(rl.KeyM) {
		currentModel = (currentModel + 1) % modelLength
		if currentModel == 0 {
			vertices = CubeVertices
			faces = CubeFaces
		} else {
			vertices = PengerVertices
			faces = PengerFaces
		}
	}
}
