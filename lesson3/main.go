package main

import "github.com/veandco/go-sdl2/sdl"

func initSDL() (*sdl.Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	return window, nil
}

func loadMedia() (*sdl.Surface, error) {
	surface, err := sdl.LoadBMP("assets/x.bmp")
	if err != nil {
		return nil, err
	}
	return surface, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	window, err := initSDL()
	must(err)
	defer window.Destroy()

	windowSurface, err := window.GetSurface()
	must(err)

	picSurface, err := loadMedia()
	must(err)
	defer picSurface.Free()

	var event sdl.Event // sdl.Event is interface{}
	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			}

			sdl.Delay(16)
		}

		picSurface.Blit(nil, windowSurface, nil)
		window.UpdateSurface()
	}

	picSurface.Blit(nil, windowSurface, nil)
	window.UpdateSurface()

	sdl.Quit()
}
