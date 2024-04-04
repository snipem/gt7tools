package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	fontPath = "./test.ttf"
	fontSize = 32
)

type history struct {
	throttle []int
	brake    []int
}

func historyToPoints(speedHistory []int, renderer *sdl.Renderer) []sdl.Point {
	var points []sdl.Point

	points = []sdl.Point{}

	// TODO get window width
	window, _ := renderer.GetWindow()
	maxwidth, maxheight := window.GetSize()

	for i := int(maxwidth); i > 0; i-- {
		if len(speedHistory) > i {
			points = append(points, sdl.Point{X: maxwidth - int32(i), Y: maxheight - int32(speedHistory[len(speedHistory)-1-i])})
		}
	}
	return points
}

//func recordData(gt7c *lib.GT7Communication, throttleHistory *[]int, brakeHistory *[]int) {
//
//	for true {
//		throttleHistory = append(throttleHistory, gt7c.LastData.Throttle)
//		brakeHistory = append(brakeHistory, gt7c.LastData.Brake)
//		time.Sleep(16 * time.Millisecond)
//	}
//}

func recordDataDemo(history *history) {

	for i := 0; i < 100; i++ {

		history.throttle = append(history.throttle, int(math.Sin(float64(i))))
		history.brake = append(history.throttle, 3+int(math.Sin(float64(i))))

		println(math.Sin(float64(i)) * 10)

		if i <= 100 {
			i = 0
		}

		time.Sleep(16 * time.Millisecond)
	}

}

func run() (err error) {
	var window *sdl.Window
	var font *ttf.Font
	var surface *sdl.Surface
	var renderer *sdl.Renderer
	var h *history
	//var points []sdl.Point
	//var rect sdl.Rect
	//var rects []sdl.Rect

	h = &history{
		throttle: []int{},
		brake:    []int{},
	}

	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return
	}
	defer sdl.Quit()

	// Create a window for us to draw the text on
	if window, err = sdl.CreateWindow("GT7 Log", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 300, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	print("Show window")
	defer window.Destroy()

	//if surface, err = window.GetSurface(); err != nil {
	//	return
	//}

	// Load the font for our text
	if font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		return
	}
	defer font.Close()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return
	}

	defer renderer.Destroy()

	// Run infinite loop until user closes the window
	running := true
	i := 0

	//gt7c := lib.NewGT7Communication("192.168.178.119")
	//go gt7c.Run()
	//go gt7c.Run()
	// Do other tasks

	//go recordData(history)
	go recordDataDemo(h)

	for running {
		for event := sdl.PollEvent(); event == nil; event = sdl.PollEvent() {

			surface.FillRect(nil, 0)

			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()
			//
			//renderer.SetDrawColor(255, 255, 255, 255)
			//renderer.DrawPoint(150, 300)
			//
			//renderer.SetDrawColor(0, 0, 255, 255)
			//renderer.DrawLine(0, 0, 200, 200)

			//points = []sdl.Point{{0, 0}, {100, 300}, {100, 300}, {200, 0}}
			//renderer.SetDrawColor(255, 255, 0, 255)
			//renderer.DrawLines(points)

			renderer.SetDrawColor(0, 255, 0, 255)

			throttleHistoryPoints := historyToPoints(h.throttle, renderer)
			if len(throttleHistoryPoints) > 0 {
				renderer.DrawLines(throttleHistoryPoints)
			}

			renderer.SetDrawColor(255, 0, 0, 255)

			brakeHistoryPoints := historyToPoints(h.brake, renderer)
			if len(brakeHistoryPoints) > 0 {
				renderer.DrawLines(brakeHistoryPoints)
			}

			//drawText(renderer, font, fmt.Sprintf("Speed:  %d", int(gt7c.LastData.CarSpeed)), 30, 30)
			//drawText(renderer, font, fmt.Sprintf("Gear:  %d", int(gt7c.LastData.CurrentGear)), 300, 30)
			renderer.Present()
			// Update the window surface with what we have drawn

			i++

			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		sdl.Delay(16)
	}

	//gt7c.Stop()
	return
}

func drawText(renderer *sdl.Renderer, font *ttf.Font, s string, x int32, y int32) {
	// Create a red text with the font
	surface, err := font.RenderUTF8Blended(s, sdl.Color{R: 255, G: 0, B: 0, A: 255})
	if err != nil {
		return
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return
	}
	defer texture.Destroy()

	renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H})
}

func main() {
	runtime.LockOSThread()
	if err := run(); err != nil {
		print("Error: %s", err)
		os.Exit(1)
	}
}
