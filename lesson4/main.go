package main

import "github.com/veandco/go-sdl2/sdl"

var window *sdl.Window

var keyPressSurfaces [5]*sdl.Surface

const (
	KEY_PRESS_DEFAULT = iota
	KEY_PRESS_UP
	KEY_PRESS_RIGHT
	KEY_PRESS_DOWN
	KEY_PRESS_LEFT
)

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

func loadBMP(path string) (*sdl.Surface, error) {
	surface, err := sdl.LoadBMP(path)
	if err != nil {
		return nil, err
	}
	return surface, nil
}

func loadMedia() {
	defaultSurface, err := loadBMP("assets/press.bmp")
	must(err)
	upSurface, err := loadBMP("assets/up.bmp")
	must(err)
	rightSurface, err := loadBMP("assets/right.bmp")
	must(err)
	downSurface, err := loadBMP("assets/down.bmp")
	must(err)
	leftSurface, err := loadBMP("assets/left.bmp")
	must(err)

	keyPressSurfaces[KEY_PRESS_DEFAULT] = defaultSurface
	keyPressSurfaces[KEY_PRESS_UP] = upSurface
	keyPressSurfaces[KEY_PRESS_RIGHT] = rightSurface
	keyPressSurfaces[KEY_PRESS_DOWN] = downSurface
	keyPressSurfaces[KEY_PRESS_LEFT] = leftSurface
}

func close() {
	for _, sur := range keyPressSurfaces {
		sur.Free()
	}

	window.Destroy()

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

	windowSurface, err := window.GetSurface()
	must(err)

	loadMedia()

	keyPressSurfaces[KEY_PRESS_DEFAULT].Blit(nil, windowSurface, nil)
	window.UpdateSurface()

	var event sdl.Event // sdl.Event is interface{}
	var quit bool
	for !quit {
		var nextSurface *sdl.Surface
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyDownEvent:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_UP:
					nextSurface = keyPressSurfaces[KEY_PRESS_UP]
				case sdl.SCANCODE_RIGHT:
					nextSurface = keyPressSurfaces[KEY_PRESS_RIGHT]
				case sdl.SCANCODE_DOWN:
					nextSurface = keyPressSurfaces[KEY_PRESS_DOWN]
				case sdl.SCANCODE_LEFT:
					nextSurface = keyPressSurfaces[KEY_PRESS_LEFT]
				}

				nextSurface.Blit(nil, windowSurface, nil)
				window.UpdateSurface()
			}
		}
	}

	close()
}

