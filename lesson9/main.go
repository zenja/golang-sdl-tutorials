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

var gTexture *sdl.Texture

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
	var err error
	gTexture, err = loadTexture("assets/viewport.png", gRenderer)
	must(err)
}

func close() {
	gTexture.Destroy()
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

		// Initialize renderer color and clear renderer
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		topLeftRect := &sdl.Rect{0, 0, SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2}
		gRenderer.SetViewport(topLeftRect)
		gRenderer.Copy(gTexture, nil, nil)

		topRightRect := &sdl.Rect{SCREEN_WIDTH / 2, 0, SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2}
		gRenderer.SetViewport(topRightRect)
		gRenderer.Copy(gTexture, nil, nil)

		bottomRect := &sdl.Rect{0, SCREEN_HEIGHT / 2, SCREEN_WIDTH, SCREEN_HEIGHT / 2}
		gRenderer.SetViewport(bottomRect)
		gRenderer.Copy(gTexture, nil, nil)

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
