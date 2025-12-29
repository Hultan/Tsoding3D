package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vector2 struct {
	X, Y float32
}

func (p Vector2) Point() {
	rl.DrawRectangle(int32(p.X)-vertexSize/2, int32(p.Y)-vertexSize/2, vertexSize, vertexSize, foregroundColor)
}

func (p Vector2) Line(p2 Vector2) {
	rl.DrawLineEx(p.ToRaylib(), p2.ToRaylib(), lineThickness, foregroundColor)
}

func (p Vector2) Screen() Vector2 {
	// -1..1 => 0..2 => 0..1 => Invert Y => 0..w/h
	p.X = (p.X + 1) / 2 * width
	p.Y = (1 - (p.Y+1)/2) * height

	return p
}

func (p Vector2) ToRaylib() rl.Vector2 {
	return rl.Vector2{X: p.X, Y: p.Y}
}

type Vector3 struct {
	X, Y, Z float32
}

func (p Vector3) TranslateZ(dz float32) Vector3 {
	return Vector3{p.X, p.Y, p.Z + dz}
}

func (p Vector3) RotateXY(angle float32) Vector3 {
	var (
		c = cos(angle)
		s = sin(angle)
	)
	return Vector3{
		X: p.X*c - p.Z*s,
		Y: p.Y,
		Z: p.X*s + p.Z*c,
	}
}

func (p Vector3) Project() Vector2 {
	return Vector2{p.X / p.Z, p.Y / p.Z}
}

// Helper functions

func cos(a float32) float32 {
	return float32(math.Cos(float64(a)))
}

func sin(a float32) float32 {
	return float32(math.Sin(float64(a)))
}
