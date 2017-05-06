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

const WALKING_ANIMATION_FRAMES = 4

var (
	gArrowTexture *lTexture
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

func (t *lTexture) setColor(r, g, b uint8) {
	t.texture.SetColorMod(r, g, b)
}

func (t *lTexture) setAlpha(alpha uint8) {
	t.texture.SetAlphaMod(alpha)
}

func (t *lTexture) setBlendMode(bm sdl.BlendMode) {
	t.texture.SetBlendMode(bm)
}

func (t *lTexture) render(x, y int32, clip *sdl.Rect) {
	renderQuad := &sdl.Rect{x, y, t.width, t.height}
	if clip != nil {
		renderQuad.W = clip.W
		renderQuad.H = clip.H
	}
	gRenderer.Copy(t.texture, clip, renderQuad)
}

func (t *lTexture) renderRotationFlip(x, y int32, clip *sdl.Rect, angle float64, center *sdl.Point, flip sdl.RendererFlip) {
	renderQuad := &sdl.Rect{x, y, t.width, t.height}
	if clip != nil {
		renderQuad.W = clip.W
		renderQuad.H = clip.H
	}
	gRenderer.CopyEx(t.texture, clip, renderQuad, angle, center, flip)
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
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)

	return window, renderer, nil
}

func loadMedia() {
	gArrowTexture = NewLTexture(gRenderer, "assets/arrow.png", nil)
}

func close() {
	gRenderer.Destroy()
	gWindow.Destroy()

	gArrowTexture.free()

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

	// Angle of rotation
	var degrees float64

	// Flip type
	var flipType sdl.RendererFlip = sdl.FLIP_NONE

	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyDownEvent:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_A:
					degrees -= 60
				case sdl.SCANCODE_D:
					degrees += 60
				case sdl.SCANCODE_Q:
					flipType = sdl.FLIP_HORIZONTAL
				case sdl.SCANCODE_W:
					flipType = sdl.FLIP_NONE
				case sdl.SCANCODE_E:
					flipType = sdl.FLIP_VERTICAL
				}
			}
		}

		// Clear screen
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		// Render arrow
		gArrowTexture.renderRotationFlip((SCREEN_WIDTH-gArrowTexture.width)/2, (SCREEN_HEIGHT-gArrowTexture.height)/2, nil, degrees, nil, flipType)

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
