package main

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/raylib"
)

func main() {
	const screenWidth = 800
	const screenHeight = 450
	raylib.InitWindow(screenWidth, screenHeight, c.Str("Raylib DEMO"))
	for !raylib.WindowShouldClose() {
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RAYWHITE)
		raylib.DrawRectangle(screenWidth/2-50, screenHeight/2-50, 100, 100, raylib.BLUE)
		raylib.EndDrawing()
	}
}
