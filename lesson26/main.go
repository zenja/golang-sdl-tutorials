package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

/* ------------------------------ global constants ------------------------------ */

const (
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
)

const (
	DOT_WIDTH  = 20
	DOT_HEIGHT = 20
)

// Maximum axis velocity of the dot
const DOT_VEL = 2

/* ------------------------------ global variables ------------------------------ */

var gWindow *sdl.Window
var gRenderer *sdl.Renderer

var gDotTexture *MyTexture

/* ------------------------------ lesson-specific types ------------------------------ */

type dot struct {
	x, y       int32
	velX, velY int32
}

func (d *dot) handleEvent(e sdl.Event) {
	switch t := e.(type) {
	case *sdl.KeyDownEvent:
		// If you're wondering why we're checking if the key repeat is 0,
		// it's because key repeat is enabled by default and if you press
		// and hold a key it will report multiple key presses.
		// That means we have to check if the key press is the first one
		// because we only care when the key was first pressed.
		if t.Repeat == 0 {
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_UP:
				d.velY -= DOT_VEL
			case sdl.SCANCODE_DOWN:
				d.velY += DOT_VEL
			case sdl.SCANCODE_LEFT:
				d.velX -= DOT_VEL
			case sdl.SCANCODE_RIGHT:
				d.velX += DOT_VEL
			}
		}
	}
}

func (d *dot) move() {
	// Move the dot left or right
	d.x += d.velX

	// If the dot went too far to the left or right
	if d.x < 0 || d.x+DOT_WIDTH > SCREEN_WIDTH {
		d.x -= d.velX
	}

	// Move the dot up or down
	d.y += d.velY

	// If the dot went too far to the up or down
	if d.y < 0 || d.y+DOT_HEIGHT > SCREEN_HEIGHT {
		d.y -= d.velY
	}
}

func (d *dot) render() {
	gDotTexture.render(d.x, d.y, nil)
}

/* ------------------------------ MyTexture ------------------------------ */

type MyTexture struct {
	// The renderer
	renderer *sdl.Renderer

	// The actual hardware texture
	texture *sdl.Texture

	// Image size
	width  int32
	height int32
}

func NewMyTexture(renderer *sdl.Renderer, path string, colorKey *sdl.Color) *MyTexture {
	t := &MyTexture{renderer: renderer}
	t.loadFromFile(path, colorKey)
	return t
}

func NewTextMyTexture(renderer *sdl.Renderer, text string, font *ttf.Font, color sdl.Color) *MyTexture {
	t := &MyTexture{renderer: renderer}
	t.loadFromRenderedText(text, color, font)
	return t
}

func (t *MyTexture) free() {
	if t.texture != nil {
		t.texture.Destroy()
	}
}

func (t *MyTexture) loadFromFile(path string, colorKey *sdl.Color) {
	// Free pre-existing texture
	t.free()

	surface, err := img.Load(path)
	must(err)

	if colorKey != nil {
		surface.SetColorKey(1, sdl.MapRGB(surface.Format, colorKey.R, colorKey.G, colorKey.B))
	}

	t.texture, err = t.renderer.CreateTextureFromSurface(surface)
	must(err)

	t.width = surface.W
	t.height = surface.H

	// Free loaded surface
	surface.Free()
}

func (t *MyTexture) loadFromRenderedText(textureText string, textureColor sdl.Color, font *ttf.Font) {
	// Free pre-existing texture
	t.free()

	surface, err := font.RenderUTF8_Solid(textureText, textureColor)
	must(err)

	t.texture, err = t.renderer.CreateTextureFromSurface(surface)
	must(err)

	t.width = surface.W
	t.height = surface.H

	// Free loaded surface
	surface.Free()
}

func (t *MyTexture) setColor(r, g, b uint8) {
	t.texture.SetColorMod(r, g, b)
}

func (t *MyTexture) setAlpha(alpha uint8) {
	t.texture.SetAlphaMod(alpha)
}

func (t *MyTexture) setBlendMode(bm sdl.BlendMode) {
	t.texture.SetBlendMode(bm)
}

func (t *MyTexture) render(x, y int32, clip *sdl.Rect) {
	renderQuad := &sdl.Rect{x, y, t.width, t.height}
	if clip != nil {
		renderQuad.W = clip.W
		renderQuad.H = clip.H
	}
	gRenderer.Copy(t.texture, clip, renderQuad)
}

func (t *MyTexture) renderRotationFlip(x, y int32, clip *sdl.Rect, angle float64, center *sdl.Point, flip sdl.RendererFlip) {
	renderQuad := &sdl.Rect{x, y, t.width, t.height}
	if clip != nil {
		renderQuad.W = clip.W
		renderQuad.H = clip.H
	}
	gRenderer.CopyEx(t.texture, clip, renderQuad, angle, center, flip)
}

/* ------------------------------ other ------------------------------ */

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

	// Init font system
	err = ttf.Init()
	if err != nil {
		return nil, nil, err
	}

	return window, renderer, nil
}

func loadMedia() {
	gDotTexture = NewMyTexture(gRenderer, "assets/dot.bmp", nil)
}

func close() {
	gRenderer.Destroy()
	gWindow.Destroy()

	gDotTexture.free()

	// Quit SDL subsystems
	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

/* ------------------------------ main ------------------------------ */

func main() {
	var err error

	gWindow, gRenderer, err = initSDL()
	must(err)

	loadMedia()

	var event sdl.Event // sdl.Event is interface{}

	var d dot

	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			}

			d.handleEvent(event)
		}

		// Clear screen
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		// Move the dot
		d.move()

		// Render the dot
		d.render()

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
