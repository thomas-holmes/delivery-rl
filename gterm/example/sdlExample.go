package main

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const tileSize = 40
const windowWidth = 640
const windowHeight = 480

func getResource(root string, asset string) string {
	return path.Join("assets", root, asset)
}

func loadTexture(file string, renderer *sdl.Renderer) (*sdl.Texture, error) {
	texture, err := img.LoadTexture(renderer, file)
	if err != nil {
		return nil, err
	}

	return texture, nil
}

func finalRenderTexture(texture *sdl.Texture, renderer *sdl.Renderer,
	destination *sdl.Rect, clip *sdl.Rect) error {

	return renderer.Copy(texture, clip, destination)
}

func renderTexture(texture *sdl.Texture, renderer *sdl.Renderer, xPos int, yPos int, clip *sdl.Rect) error {
	_, _, width, height, err := texture.Query()
	if err != nil {
		return err
	}

	dest := sdl.Rect{H: height, W: width, X: int32(xPos), Y: int32(yPos)}

	return finalRenderTexture(texture, renderer, &dest, clip)
}

func renderTextureScaled(texture *sdl.Texture, renderer *sdl.Renderer, x int, y int, w int, h int) error {
	dest := sdl.Rect{H: int32(h), W: int32(w), X: int32(x), Y: int32(y)}
	err := renderer.Copy(texture, nil, &dest)
	if err != nil {
		return err
	}
	return nil
}

func renderText(renderer *sdl.Renderer, message string, fontFile string, color sdl.Color, fontSize int) (*sdl.Texture, error) {

	font, err := ttf.OpenFont(fontFile, fontSize)
	if err != nil {
		return nil, err
	}
	defer font.Close()

	surface, err := font.RenderUTF8_Blended(message, color)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	return texture, err
}

func setSdlLogger() {
	/*
		LOG_PRIORITY_VERBOSE
		LOG_PRIORITY_DEBUG
		LOG_PRIORITY_INFO
		LOG_PRIORITY_WARN
		LOG_PRIORITY_ERROR
		LOG_PRIORITY_CRITICAL
	*/
	sdl.LogSetOutputFunction(func(data interface{}, cat int, pri sdl.LogPriority, message string) {
		priArray := [6]string{"VERBOSE", "DEBUG", "INFO", "WARN", "ERROR", "CRITICAL"}
		log.Println("[SDL]", fmt.Sprintf("[%v]", priArray[pri-1]), message)
	}, nil)
}

func tileBackground(background *sdl.Texture, renderer *sdl.Renderer, width int, height int, tileSize int) error {
	xTiles := width / tileSize
	yTiles := height / tileSize

	for tile := 0; tile < xTiles*yTiles; tile++ {
		x := tile % xTiles
		y := tile / xTiles

		err := renderTextureScaled(background, renderer, x*tileSize, y*tileSize, tileSize, tileSize)
		if err != nil {
			return err
		}
	}
	return nil
}

func initSdl() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	ttf.Init()
}

func createWindow() *sdl.Window {
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	return window
}

func createRenderer(window *sdl.Window) *sdl.Renderer {
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}
	return renderer
}

func main() {
	initSdl()
	defer sdl.Quit()

	window := createWindow()
	defer window.Destroy()

	renderer := createRenderer(window)
	defer renderer.Destroy()

	background, err := loadTexture(getResource("img", "background.png"), renderer)
	if err != nil {
		panic(err)
	}
	defer background.Destroy()

	clipMap := [4]sdl.Rect{}
	clipWidth := 100
	clipHeight := 100
	for i := 0; i < 4; i++ {
		clip := sdl.Rect{
			X: int32(i / 2 * clipWidth),
			Y: int32(i % 2 * clipHeight),
			W: int32(clipWidth),
			H: int32(clipHeight),
		}
		clipMap[i] = clip
	}

	foreground, err := loadTexture(getResource("img", "image.png"), renderer)
	if err != nil {
		panic(err)
	}
	defer foreground.Destroy()

	quit := false
	frameCount := 0
	fps := 0
	clipIndex := 0

	go startFpsCounter(&frameCount, &fps)
	for quit == false {
		renderer.Clear()
		event := sdl.PollEvent()
		if event != nil {
			// log.Println(fmt.Sprintf("got event %#v", event))
			switch e := event.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyDownEvent:
				switch e.Keysym.Sym {
				case sdl.K_1:
					clipIndex = 0
				case sdl.K_2:
					clipIndex = 1
				case sdl.K_3:
					clipIndex = 2
				case sdl.K_4:
					clipIndex = 3
				}
			case *sdl.MouseButtonEvent:
				quit = true
			}
		}
		frameCount++

		tileBackground(background, renderer, windowWidth, windowHeight, tileSize)
		renderTexture(foreground, renderer, 50, 50, &clipMap[clipIndex])
		renderFps(renderer, fps)

		renderer.Present()
		window.UpdateSurface()
	}
}

func renderFps(renderer *sdl.Renderer, fps int) {
	texture, err := renderText(renderer, fmt.Sprint(fps), getResource("font", "sample.ttf"), sdl.Color{R: 255, G: 255, B: 255}, 48)
	if err != nil {
		panic(err)
	}

	_, _, width, height, err := texture.Query()
	if err != nil {
		panic(err)
	}
	dest := sdl.Rect{X: windowWidth - width, Y: 0, W: width, H: height}
	renderer.Copy(texture, nil, &dest)
}

func startFpsCounter(i *int, fps *int) {
	timer := time.Tick(1 * time.Second)
	go func() {
		for {
			<-timer
			log.Println(*i)
			*fps = *i
			*i = 0
		}
	}()
}

func init() {
	setSdlLogger()
}