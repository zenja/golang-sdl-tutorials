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

var gLTexture *lTexture

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

func (t *lTexture) setColor(r, g, b uint8) {
	t.texture.SetColorMod(r, g, b)
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
	gLTexture = NewLTexture(gRenderer, "assets/colors.png", nil)
}

func close() {
	gRenderer.Destroy()
	gWindow.Destroy()

	gLTexture.free()

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
	var r uint8 = 255
	var g uint8 = 255
	var b uint8 = 255
	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyDownEvent:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_Q:
					r += 32
				case sdl.SCANCODE_W:
					g += 32
				case sdl.SCANCODE_E:
					b += 32
				case sdl.SCANCODE_A:
					r -= 32
				case sdl.SCANCODE_S:
					g -= 32
				case sdl.SCANCODE_D:
					b -= 32
				}
			}
		}

		// Clear screen
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		//Modulate and render texture
		gLTexture.setColor(r, g, b)
		gLTexture.render(0, 0, nil)

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
