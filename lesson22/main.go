package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"strconv"
)

/* ------------------------------ global constants ------------------------------ */

const (
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
)

/* ------------------------------ global variables ------------------------------ */

var gWindow *sdl.Window
var gRenderer *sdl.Renderer

var gFont *ttf.Font
var gPromptTextTexture *MyTexture

/* ------------------------------ lesson-specific types ------------------------------ */

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
	var err error
	gFont, err = ttf.OpenFont("assets/lazy.ttf", 28)
	must(err)
	gPromptTextTexture = NewTextMyTexture(
		gRenderer, "Press Enter to Reset Start Time.", gFont, sdl.Color{0, 0, 0, 255})
}

func close() {
	gRenderer.Destroy()
	gWindow.Destroy()

	gPromptTextTexture.free()

	gFont.Close()

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

	var startTime uint32

	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyDownEvent:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_RETURN:
					startTime = sdl.GetTicks()
				}
			}

		}

		// Clear screen
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		// Render prompt texture
		gPromptTextTexture.render((SCREEN_WIDTH-gPromptTextTexture.width)/2, 0, nil)

		// Create and render time texture
		timeText := "Milliseconds since start time: " + strconv.Itoa(int(sdl.GetTicks()-startTime))
		timeTextTexture := NewTextMyTexture(gRenderer, timeText, gFont, sdl.Color{0, 0, 0, 255})
		timeTextTexture.render(
			(SCREEN_WIDTH-gPromptTextTexture.width)/2, (SCREEN_HEIGHT-gPromptTextTexture.height)/2, nil)
		timeTextTexture.free()

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
