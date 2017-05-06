package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

const (
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
)

var gWindow *sdl.Window
var gRenderer *sdl.Renderer

func initSDL() (*sdl.Window, *sdl.Renderer, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, nil, err
	}

	// Create window
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, nil, err
	}

	// Create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	return window, renderer, nil
}

func loadTexture(path string, renderer *sdl.Renderer) (*sdl.Texture, error) {
	surface, err := img.Load(path)
	if err != nil {
		return nil, err
	}
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	// Free loaded surface
	surface.Free()

	return texture, nil
}

func loadMedia() {
}

func close() {
	gRenderer.Destroy()
	gWindow.Destroy()

	// Quit SDL subsystems
	img.Quit()
	sdl.Quit()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	gWindow, gRenderer, err = initSDL()
	must(err)

	loadMedia()

	var event sdl.Event // sdl.Event is interface{}
	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			}
		}

		// Initialize renderer color
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		// Render red filled quad
		fillRect := &sdl.Rect{
			SCREEN_WIDTH / 4,
			SCREEN_HEIGHT / 4,
			SCREEN_WIDTH / 2,
			SCREEN_HEIGHT / 2,
		}
		gRenderer.SetDrawColor(255, 0, 0, 255)
		gRenderer.FillRect(fillRect)

		// Render green outlined quad
		outlineRect := &sdl.Rect{
			SCREEN_WIDTH / 6,
			SCREEN_HEIGHT / 6,
			SCREEN_WIDTH * 2 / 3,
			SCREEN_HEIGHT * 2 / 3,
		}
		gRenderer.SetDrawColor(0, 255, 0, 255)
		gRenderer.DrawRect(outlineRect)

		// Draw blue horizontal line
		gRenderer.SetDrawColor(0, 0, 255, 255)
		gRenderer.DrawLine(0, SCREEN_HEIGHT/2, SCREEN_WIDTH, SCREEN_HEIGHT/2)

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
