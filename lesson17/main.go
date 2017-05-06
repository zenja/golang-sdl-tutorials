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
	BUTTON_WIDTH  = 300
	BUTTON_HEIGHT = 200
)

const BUTTON_NUM = 4
const BUTTON_SPRITE_TOTAL = 4

const (
	BUTTON_SPRITE_MOUSE_OUT = iota
	BUTTON_SPRITE_MOUSE_OVER_MOTION
	BUTTON_SPRITE_MOUSE_DOWN
	BUTTON_SPRITE_MOUSE_UP
)

/* ------------------------------ global variables ------------------------------ */

var gWindow *sdl.Window
var gRenderer *sdl.Renderer

var gButtons [BUTTON_NUM]*button

/* ------------------------------ lesson-specific types ------------------------------ */

type buttonSpriteIndex int

type button struct {
	// Top left position
	position sdl.Point

	spriteSheetTexture *MyTexture

	spriteIndex buttonSpriteIndex

	spriteClips [BUTTON_SPRITE_TOTAL]sdl.Rect
}

func NewButton(position sdl.Point, spriteSheetTexture *MyTexture) *button {
	b := &button{
		position:           position,
		spriteSheetTexture: spriteSheetTexture,
	}
	for i := 0; i < BUTTON_SPRITE_TOTAL; i++ {
		b.spriteClips[i] = sdl.Rect{0, int32(i * 200), BUTTON_WIDTH, BUTTON_HEIGHT}
	}
	return b
}

func (b *button) setPosition(x, y int32) {
	b.position.X = x
	b.position.Y = y
}

func (b *button) render() {
	b.spriteSheetTexture.render(b.position.X, b.position.Y, &b.spriteClips[b.spriteIndex])
}

func (b *button) handleEvent(e sdl.Event) {
	var mouseEventType uint32
	switch t := e.(type) {
	case *sdl.MouseMotionEvent:
		mouseEventType = t.Type
	case *sdl.MouseButtonEvent:
		mouseEventType = t.Type
	default:
		return
	}

	var inside bool = true
	x, y, _ := sdl.GetMouseState()
	if int32(x) < b.position.X || int32(x) > b.position.X+BUTTON_WIDTH {
		inside = false
	}
	if int32(y) < b.position.Y || int32(y) > b.position.Y+BUTTON_HEIGHT {
		inside = false
	}

	if !inside {
		b.spriteIndex = BUTTON_SPRITE_MOUSE_OUT
	} else {
		if mouseEventType == sdl.MOUSEMOTION {
			b.spriteIndex = BUTTON_SPRITE_MOUSE_OVER_MOTION
		}
		if mouseEventType == sdl.MOUSEBUTTONDOWN {
			b.spriteIndex = BUTTON_SPRITE_MOUSE_DOWN
		}
		if mouseEventType == sdl.MOUSEBUTTONUP {
			b.spriteIndex = BUTTON_SPRITE_MOUSE_UP
		}
	}
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
	ttf.Init()

	return window, renderer, nil
}

func loadMedia() {
	spriteSheetTexture := NewMyTexture(gRenderer, "assets/button.png", nil)
	for i := 0; i < BUTTON_NUM; i++ {
		gButtons[i] = NewButton(sdl.Point{}, spriteSheetTexture)
	}
	gButtons[0].setPosition(0, 0)
	gButtons[1].setPosition(SCREEN_WIDTH-BUTTON_WIDTH, 0)
	gButtons[2].setPosition(0, SCREEN_HEIGHT-BUTTON_HEIGHT)
	gButtons[3].setPosition(SCREEN_WIDTH-BUTTON_WIDTH, SCREEN_HEIGHT-BUTTON_HEIGHT)
}

func close() {
	gRenderer.Destroy()
	gWindow.Destroy()

	gButtons[0].spriteSheetTexture.free()

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

	var quit bool
	for !quit {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit = true
			}

			for i := 0; i < BUTTON_NUM; i++ {
				gButtons[i].handleEvent(event)
			}
		}

		// Clear screen
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		// Render text
		for i := 0; i < BUTTON_NUM; i++ {
			gButtons[i].render()
		}

		// Update screen
		gRenderer.Present()

		sdl.Delay(16)
	}

	close()
}
