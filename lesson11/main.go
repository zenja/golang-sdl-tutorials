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

var (
	gSpriteClips         [4]*sdl.Rect
	gSpriteSheetLTexture *lTexture
)

type rgb struct {
	r, g, b uint8
}

type lTexture struct {
	// The renderer
	renderer *sdl.Renderer

	// The actual hardware texture
	texture *sdl.Texture

	// Image size
	width  int32
	height int32
}

func NewLTexture(renderer *sdl.Renderer, path string, colorKey *rgb) *lTexture {
	t := &lTexture{renderer: renderer}
	t.loadFromFile(path, colorKey)
	return t
}

func (t *lTexture) free() {
	if t.texture != nil {
		t.texture.Destroy()
	}
}

func (t *lTexture) loadFromFile(path string, colorKey *rgb) {
	// Free pre-existing texture
	t.free()

	surface, err := img.Load(path)
	must(err)

	if colorKey != nil {
		surface.SetColorKey(1, sdl.MapRGB(surface.Format, colorKey.r, colorKey.g, colorKey.b))
	}

	t.texture, err = t.renderer.CreateTextureFromSurface(surface)
	must(err)

	t.width = surface.W
	t.height = surface.H

	// Free loaded surface
	surface.Free()
}

func (t *lTexture) render(x, y int32, clip *sdl.Rect) {
	renderQuad := &sdl.Rect{x, y, t.width, t.height}
	if clip != nil {
		renderQuad.W = clip.W
		renderQuad.H = clip.H
	}
	gRenderer.Copy(t.texture, clip, renderQuad)
}

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

func loadMedia() {
	gSpriteSheetLTexture = NewLTexture(gRenderer, "assets/dots.png", &rgb{0, 255, 255})

	// Set top left sprite
	gSpriteClips[0] = &sdl.Rect{0, 0, 100, 100}

	// Set top right sprite
	gSpriteClips[1] = &sdl.Rect{100, 0, 100, 100}

	// Set bottom left sprite
	gSpriteClips[2] = &sdl.Rect{0, 100, 100, 100}

	// Set bottom right sprite
	gSpriteClips[3] = &sdl.Rect{100, 100, 100, 100}
}

func close() {
	gRenderer.Destroy()
	gWindow.Destroy()

	gSpriteSheetLTexture.free()

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

		//Render top left sprite
		gSpriteSheetLTexture.render(0, 0, gSpriteClips[0])

		//Render top right sprite
		gSpriteSheetLTexture.render(SCREEN_WIDTH-gSpriteClips[1].W, 0, gSpriteClips[1])

		//Render bottom left sprite
		gSpriteSheetLTexture.render(0, SCREEN_HEIGHT-gSpriteClips[2].H, gSpriteClips[2])

		//Render bottom right sprite
		gSpriteSheetLTexture.render(SCREEN_WIDTH-gSpriteClips[3].W, SCREEN_HEIGHT-gSpriteClips[3].H, gSpriteClips[3])

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
