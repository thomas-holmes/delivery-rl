package gterm

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Window represents the base window object
type Window struct {
	Columns         int
	Rows            int
	FontSize        int
	tileHeightPixel int
	tileWidthPixel  int
	heightPixel     int
	widthPixel      int
	fontPath        string
	font            *ttf.Font
	SdlWindow       *sdl.Window
	SdlRenderer     *sdl.Renderer
	backgroundColor sdl.Color
	cells           [][]renderItem
	fps             fpsCounter
	fpsLimit        int
	drawInterval    uint32
}

type renderItem struct {
	FColor  sdl.Color
	Glyph   string
	Texture *sdl.Texture
}

func (renderItem *renderItem) Destroy() {
	if renderItem.Texture != nil {
		renderItem.Texture.Destroy()
		renderItem.Texture = nil
	}
}

// NewWindow constructs a window
func NewWindow(columns int, rows int, fontPath string, fontSize int, fpsLimit int) *Window {

	numCells := columns * rows
	cells := make([][]renderItem, numCells, numCells)
	drawInterval := uint32(0)
	if fpsLimit != 0 {
		drawInterval = uint32(1000 / fpsLimit)
	}

	window := &Window{
		Columns:      columns,
		Rows:         rows,
		FontSize:     fontSize,
		fontPath:     fontPath,
		cells:        cells,
		fpsLimit:     fpsLimit,
		drawInterval: drawInterval,
	}

	return window
}

func computeCellSize(font *ttf.Font) (width int, height int, err error) {
	atGlyph, err := font.RenderUTF8_Blended("@", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return 0, 0, err
	}
	return int(atGlyph.W), int(atGlyph.H), nil
}

func (window *Window) SetTitle(title string) {
	window.SdlWindow.SetTitle(title)
}

// Init initialized the window for drawing
func (window *Window) Init() error {
	err := sdl.Init(sdl.INIT_EVERYTHING) // not sure where to do this
	if err != nil {
		return err
	}
	err = ttf.Init()
	if err != nil {
		return nil
	}
	openedFont, err := ttf.OpenFont(window.fontPath, window.FontSize)
	if err != nil {
		return err
	}

	window.font = openedFont
	tileWidth, tileHeight, err := computeCellSize(window.font)
	if err != nil {
		return err
	}

	window.tileWidthPixel = tileWidth
	window.tileHeightPixel = tileHeight
	window.heightPixel = tileHeight * window.Rows
	window.widthPixel = tileWidth * window.Columns

	log.Printf("Creating window w:%v, h:%v", window.widthPixel, window.heightPixel)
	sdlWindow, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, window.widthPixel, window.heightPixel, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}

	sdlRenderer, err := sdl.CreateRenderer(sdlWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}

	err = sdlRenderer.SetDrawColor(10, 10, 25, 255)
	if err != nil {
		log.Fatalln("Could not set render color", err)
	}

	window.SdlWindow = sdlWindow
	window.SdlRenderer = sdlRenderer

	window.fps = newFpsCounter()

	return nil
}

func (window *Window) SetBackgroundColor(color sdl.Color) {
	window.backgroundColor = color
}

func (window *Window) cellIndex(col int, row int) (int, error) {
	if col >= window.Columns || col < 0 || row >= window.Rows || row < 0 {
		return 0, fmt.Errorf("Requested invalid position (%v,%v) on board of dimensions %vx%v", col, row, window.Columns, window.Rows)
	}
	return col + window.Columns*row, nil
}

func (window *Window) renderCell(col int, row int) error {
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}

	renderItems := window.cells[index]

	destinationRect := sdl.Rect{
		X: int32(col * window.tileWidthPixel),
		Y: int32(row * window.tileHeightPixel),
		W: int32(window.tileWidthPixel),
		H: int32(window.tileHeightPixel),
	}

	for index := range renderItems {
		renderItem := &renderItems[index]
		// surface, err := window.font.RenderUTF8_Blended(renderItem.Glyph, renderItem.FColor)
		if renderItem.Texture == nil {
			surface, err := window.font.RenderUTF8_Solid(renderItem.Glyph, renderItem.FColor)
			if err != nil {
				return err
			}
			defer surface.Free()

			texture, err := window.SdlRenderer.CreateTextureFromSurface(surface)
			if err != nil {
				return err
			}
			renderItem.Texture = texture
		}
		_, _, width, height, err := renderItem.Texture.Query()
		if err != nil {
			return err
		}

		destinationRect.W = width
		destinationRect.H = height

		// log.Println(destinationRect)
		err = window.SdlRenderer.Copy(renderItem.Texture, nil, &destinationRect)
		if err != nil {
			return err
		}
	}

	return nil
}

func (window *Window) renderCells() error {
	for col := 0; col < window.Columns; col++ {
		for row := 0; row < window.Rows; row++ {
			if err := window.renderCell(col, row); err != nil {
				return err
			}
		}
	}
	return nil
}

func (window *Window) AddToCell(col int, row int, glyph string, fColor sdl.Color) error {
	renderItem := renderItem{Glyph: glyph, FColor: fColor}
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}
	window.cells[index] = append(window.cells[index], renderItem)

	return nil
}

func (window *Window) ClearCell(col int, row int) error {
	index, err := window.cellIndex(col, row)
	if err != nil {
		return err
	}

	for _, item := range window.cells[index] {
		item.Destroy()
	}
	window.cells[index] = window.cells[index][:0]

	return nil
}

func (window *Window) ClearWindow() {
	window.cells = make([][]renderItem, window.Columns*window.Rows, window.Columns*window.Rows)
}

func (window *Window) ShouldRenderFps(shouldRender bool) {
	window.fps.shouldRender(shouldRender)
}

// Render updates the display based on new information since last Render
type fpsCounter struct {
	renderFps     bool
	framesElapsed int
	currentFps    int
	lastTicks     uint32
	font          *ttf.Font
	color         sdl.Color
}

func newFpsCounter() fpsCounter {
	font, err := ttf.OpenFont(path.Join("assets/font/FiraMono-Regular.ttf"), 16)
	if err != nil {
		log.Println("Failed to open FPS font", err)
	}
	return fpsCounter{
		lastTicks: sdl.GetTicks(),
		color:     sdl.Color{R: 0, G: 255, B: 0, A: 255},
		font:      font,
	}
}

func (fps *fpsCounter) shouldRender(shouldRender bool) {
	fps.renderFps = shouldRender
}

func (fps *fpsCounter) MaybeRender(window *Window) {
	fps.framesElapsed++
	now := sdl.GetTicks()

	if fps.lastTicks < (now - 1000) {
		fps.lastTicks = now
		fps.currentFps = fps.framesElapsed
		fps.framesElapsed = 0
	}

	if fps.renderFps {
		surface, err := fps.font.RenderUTF8_Solid(strconv.Itoa(fps.currentFps), fps.color)
		if err != nil {
			log.Println("Failed to render FPS", err)
		}
		defer surface.Free()
		texture, err := window.SdlRenderer.CreateTextureFromSurface(surface)
		if err != nil {
			log.Println("Failed to texturify FPS", err)
		}
		defer texture.Destroy()

		_, _, width, height, err := texture.Query()
		destination := sdl.Rect{X: int32(window.widthPixel) - width, Y: int32(0), W: width, H: height}

		if err := window.SdlRenderer.Copy(texture, nil, &destination); err != nil {
			log.Println("Couldn't copy FPS to renderer", err)
		}
	}
}

var lastDraw = sdl.GetTicks()

func min(a uint32, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
func (window *Window) Render() {
	err := window.SdlRenderer.SetDrawColor(window.backgroundColor.R, window.backgroundColor.G, window.backgroundColor.B, window.backgroundColor.A)
	if err != nil {
		log.Fatal(err)
	}
	window.SdlRenderer.Clear()

	window.renderCells()

	window.fps.MaybeRender(window) // Ok this is dumb

	window.SdlRenderer.Present()

	if window.fpsLimit > 0 {
		nowTicks := sdl.GetTicks()
		if lastDraw < (nowTicks - window.drawInterval) {
			delay := window.drawInterval - min(0, (nowTicks-lastDraw))
			sdl.Delay(delay)
		}
	}
}