package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
)

var stretchedSurface *sdl.Surface
var window *sdl.Window

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

func loadBMP(path string) (*sdl.Surface, error) {
	surface, err := sdl.LoadBMP(path)
	if err != nil {
		return nil, err
	}
	return surface, nil
}

func loadMedia() {
	var err error
	stretchedSurface, err = loadBMP("assets/stretch.bmp")
	must(err)
}

func close() {
	stretchedSurface.Free()
	window.Destroy()
	sdl.Quit()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	window, err := initSDL()
	must(err)

	windowSurface, err := window.GetSurface()
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

		stretchRect := sdl.Rect{X: 0, Y: 0, W: SCREEN_WIDTH, H: SCREEN_HEIGHT}
		stretchedSurface.BlitScaled(nil, windowSurface, &stretchRect)
		window.UpdateSurface()

		sdl.Delay(16)
	}

	close()
}
