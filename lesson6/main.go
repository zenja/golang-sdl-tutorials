package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

const (
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
)

var window *sdl.Window
var windowSurface *sdl.Surface

var loadedSurface *sdl.Surface

func initSDL() (*sdl.Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	var err error
	window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	return window, nil
}

func loadSurface(path string) (*sdl.Surface, error) {
	var err error

	surface, err := img.Load(path)
	if err != nil {
		return nil, err
	}

	//Convert surface to screen format
	optimizedSurface, err := surface.Convert(windowSurface.Format, 0)

	//Get rid of old loaded surface
	surface.Free()

	return optimizedSurface, nil
}

func loadMedia() {
	var err error
	loadedSurface, err = loadSurface("assets/loaded.png")
	must(err)
}

func close() {
	loadedSurface.Free()
	window.Destroy()

	// quit SDL subsystems
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

	window, err = initSDL()
	must(err)

	windowSurface, err = window.GetSurface()
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

		loadedSurface.Blit(nil, windowSurface, nil)
		window.UpdateSurface()

		sdl.Delay(16)
	}

	close()
}
