package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/text"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/snipem/gt7tools/lib"
	"github.com/snipem/gt7tools/lib/dump"
	"log"
	"math"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
)

var gt7c = &gt7.GT7Communication{}

const show_n_values = 500

var showBrake bool
var showThrottle bool
var showGear bool
var signalRisingTrailbreak bool
var dumpFile string

var showTrainingBars bool

func createArrayWithValues(value int, max_values int) []float64 {
	var s []float64
	for i := 0; i < max_values; i++ {
		s = append(s, float64(value))
	}
	return s
}

// sineInputs generates values from -1 to 1 for display on the line chart.
func sineInputs() []float64 {
	var res []float64

	for i := 0; i < 200; i++ {
		v := math.Sin(float64(i) / 100 * math.Pi)
		res = append(res, v)
	}
	return res
}

func convertIntSliceToFloatSlice(intSlice []int) []float64 {
	floatSlice := make([]float64, len(intSlice))
	for i, num := range intSlice {
		floatSlice[i] = float64(num)
	}
	return floatSlice
}

func takeLastN(slice []int, n int) []int {
	if len(slice) <= n {
		return slice
	}
	startIndex := len(slice) - n
	return slice[startIndex:]
}

// playLineChart continuously adds values to the LineChart, once every delay.
// Exits when the context expires.
func playLineChart(ctx context.Context, lc *linechart.LineChart, history *lib.History, delay time.Duration) {
	inputs := sineInputs()
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for i := 0; ; {
		select {
		case <-ticker.C:

			if gt7c.LastData.IsPaused {
				// Skip recalculation if paused, makes the graph pause
				continue
			}

			updateSeries(lc, showThrottle, "throttle", history.Throttle, show_n_values, cell.FgColor(cell.ColorNumber(64)))

			updateSeries(lc, true, "position", mapToInt(history.Position_Z, 0, 100), show_n_values, cell.FgColor(cell.ColorMagenta))

			if showGear {
				// TODO get this from telemetry
				maxGear := 8
				scale := 100

				i = (i + 1) % len(inputs)
				if err := lc.Series("gear", convertIntSliceToFloatSlice(mapGearToScale(maxGear, scale, takeLastN(history.Gear, show_n_values))),
					linechart.SeriesCellOpts(cell.FgColor(cell.ColorGray)),
				); err != nil {
					panic(err)
				}

			} else {
				lc.Series("gear", createArrayWithValues(0, show_n_values), linechart.SeriesCellOpts(cell.FgColor(cell.ColorGray)))
			}

			if showBrake {
				brakeColor := cell.ColorNumber(160)

				if signalRisingTrailbreak && breakingIncreasing(history) {
					// Braking increasing after reaching peak
					brakeColor = cell.ColorBlue
				}

				if err := lc.Series("braking", convertIntSliceToFloatSlice(takeLastN(history.Brake, show_n_values)),
					linechart.SeriesCellOpts(cell.FgColor(brakeColor)),
				); err != nil {
					panic(err)
				}
			} else {
				lc.Series("braking", createArrayWithValues(0, show_n_values), linechart.SeriesCellOpts(cell.FgColor(cell.ColorWhite)))
			}

			// Static bars
			if showTrainingBars {
				trainingColor := cell.FgColor(cell.ColorWhite)

				if signalRisingTrailbreak && breakingIncreasing(history) {
					// Braking increasing after reaching peak
					trainingColor = cell.BgColor(cell.ColorBlue)
					warnBar
				}

				//if history.Brake[len(history.Brake)-1] == 100 {
				//	trainingColor = cell.BgColor(cell.ColorRed)
				//else if gt7c.LastData.IsTCSEngaged {
				//	trainingColor = cell.BgColor(cell.ColorBlue)
				//}

				trainingBars := []int{25, 50, 75, 99}

				for _, trainingBar := range trainingBars {

					if err := lc.Series(fmt.Sprintf("%d", trainingBar), createArrayWithValues(trainingBar, show_n_values),
						linechart.SeriesCellOpts(trainingColor),
					); err != nil {
						panic(err)
					}

				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func mapToInt(z []float32, min int, max int) (out []int) {
	out = make([]int, len(z))
	for i, f := range z {
		out[i] = int(f*float32(max-min)) + min
	}
	return out
}

func updateSeries(lc *linechart.LineChart, showSeries bool, label string, historyValues []int, showNValues int, color cell.Option) {
	if showSeries {
		if err := lc.Series(label, convertIntSliceToFloatSlice(takeLastN(historyValues, showNValues)),
			linechart.SeriesCellOpts(color),
		); err != nil {
			panic(err)
		}
	} else {
		lc.Series(label, createArrayWithValues(0, showNValues), linechart.SeriesCellOpts(color))
	}
}

func breakingIncreasing(history *lib.History) bool {
	return history.Brake[len(history.Brake)-1] > history.Brake[len(history.Brake)-2] &&
		!straightIncreaseFromZeroBraking(history.Brake)
}

func mapGearToScale(maxGear int, scale int, originalGear []int) (mappedGears []int) {
	multiplier := scale / maxGear

	for i := 0; i < len(originalGear); i++ {
		mappedGears = append(mappedGears, originalGear[i]*multiplier)
	}
	return mappedGears
}

func straightIncreaseFromZeroBraking(brake []int) bool {

	// Ignore all values prior to last full brake
	valuesSinceLastFullbrake := []int{}
	for i := len(brake) - 1; i > 0; i-- {
		valuesSinceLastFullbrake = append([]int{brake[i]}, valuesSinceLastFullbrake...)
		if brake[i] == 0 {
			break
		}
	}

	if len(valuesSinceLastFullbrake) == 0 {
		return false
	}

	for i := len(valuesSinceLastFullbrake) - 1; i > 0; i-- {
		if valuesSinceLastFullbrake[i] < valuesSinceLastFullbrake[i-1] {
			return false
		}
		if valuesSinceLastFullbrake[i] == 0 {
			break
		}
	}

	return true
}

// playBarChart continuously changes the displayed values on the bar chart once every delay.
// Exits when the context expires.
func playBarChart(ctx context.Context, bc *barchart.BarChart, delay time.Duration) {
	const max = 100

	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var values []int
			//for i := 0; i < 2; i++ {
			//	values = append(values, int(rand.Int31n(max+1)))
			//}
			values = []int{
				int(gt7c.LastData.Throttle),
				int(gt7c.LastData.Brake),
			}

			if err := bc.Values(values, max); err != nil {
				panic(err)
			}

		case <-ctx.Done():
			return
		}
	}
}

var warnBar *Text

func Run() {

	if dumpFile != "" {
		gt7d, err := dump.NewGT7Dump(dumpFile, gt7c)
		if err != nil {
			log.Fatal(err)
		}
		go gt7d.Run()
	} else {
		gt7c = gt7.NewGT7Communication("255.255.255.255")
		go gt7c.Run()
	}

	history := &lib.History{
		Throttle: make([]int, show_n_values),
		Brake:    make([]int, show_n_values),
		Gear:     make([]int, show_n_values),
	}

	go lib.UpdateHistory(gt7c, history)

	//const redrawInterval = 16 * time.Millisecond // 60 FPS
	const redrawInterval = 32 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
		linechart.YAxisCustomScale(0, 100),
		linechart.YAxisFormattedValues(func(value float64) string {
			return fmt.Sprintf("%d", int(value))
		}),
	)

	warnChart, err := barchart.New()

	bc, err := barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorGreen,
			cell.ColorRed,
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorBlack,
			cell.ColorBlack,
		}),
		barchart.ShowValues(),
		barchart.BarWidth(8),
		barchart.Labels([]string{
			"Throttle",
			"Brake",
		}),
	)

	if err != nil {
		panic(err)
	}
	go playLineChart(ctx, lc, history, redrawInterval)
	go playBarChart(ctx, bc, redrawInterval)

	warnBar, err = text.New()
	if err != nil {
		panic(err)
	}

	tbx, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer tbx.Close()

	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(
			89,
			grid.ColWidthPerc(85, grid.Widget(lc)),
			grid.ColWidthPerc(15, grid.Widget(bc)),
		),
		grid.RowHeightPerc(
			10,
			grid.ColWidthPerc(99, grid.Widget()),
		),
	)
	gridOpts, err := builder.Build()
	if err != nil {
		panic(err)
	}

	cont, err := container.New(tbx, gridOpts...)
	if err != nil {
		panic(err)
	}

	//c, err := container.New(
	//	t,
	//	container.Border(linestyle.Light),
	//	container.BorderTitle("PRESS Q TO QUIT"),
	//	container.PlaceWidget(lc),
	//)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}

		if k.Key == 'p' {
			showTrainingBars = !showTrainingBars
		}

		if k.Key == 't' {
			showThrottle = !showThrottle
		}

		if k.Key == 'b' {
			showBrake = !showBrake
		}

		if k.Key == 'g' {
			showGear = !showGear
		}
	}

	if err := termdash.Run(ctx, tbx, cont, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}
}

func main() {

	flag.BoolVar(&showTrainingBars, "show-training-bars", true, "Show training bars")
	flag.BoolVar(&showBrake, "show-brake", true, "Show brake")
	flag.BoolVar(&signalRisingTrailbreak, "signal-rising-trailbreak", true, "Signal rising trailbreak")
	flag.BoolVar(&showThrottle, "show-throttle", false, "Show throttle")
	flag.BoolVar(&showGear, "show-gear", false, "Show gear mapped to scale")
	flag.StringVar(&dumpFile, "dump-file", "", "Use dump file for dump values instead of telemetry. Leave blank for real telemetry.")

	flag.Parse()

	Run()
}
