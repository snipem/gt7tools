package main

import (
	"context"
	"fmt"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/widgets/barchart"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/snipem/gt7-utils/lib"
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

const showTrainingBars = true

// playLineChart continuously adds values to the LineChart, once every delay.
// Exits when the context expires.
func playLineChart(ctx context.Context, lc *linechart.LineChart, history *lib.History, delay time.Duration) {
	inputs := sineInputs()
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for i := 0; ; {
		select {
		case <-ticker.C:
			i = (i + 1) % len(inputs)
			if err := lc.Series("throttle", convertIntSliceToFloatSlice(takeLastN(history.Throttle, show_n_values)),
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(64))),
			); err != nil {
				panic(err)
			}

			if err := lc.Series("braking", convertIntSliceToFloatSlice(takeLastN(history.Brake, show_n_values)),
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(160))),
			); err != nil {
				panic(err)
			}

			// Static bars
			if showTrainingBars {
				trainingBars := []int{25, 50, 75, 99}

				for _, trainingBar := range trainingBars {

					if err := lc.Series(fmt.Sprintf("%d", trainingBar), createArrayWithValues(trainingBar, show_n_values),
						linechart.SeriesCellOpts(cell.FgColor(cell.ColorWhite)),
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

func main() {
	//t, err := tcell.New()
	//if err != nil {
	//	panic(err)
	//}
	//defer t.Close()

	gt7c = gt7.NewGT7Communication("255.255.255.255")
	go gt7c.Run()

	history := &lib.History{
		Throttle: make([]int, show_n_values),
		Brake:    make([]int, show_n_values),
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

	tbx, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer tbx.Close()

	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(
			99,
			grid.ColWidthPerc(85, grid.Widget(lc)),
			grid.ColWidthPerc(15, grid.Widget(bc)),
		),
		//grid.RowHeightPerc(
		//	50,
		//	grid.ColWidthPerc(50, grid.Widget(lc)),
		//	grid.ColWidthPerc(50, grid.Widget(lc)),
		//),
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
	}

	if err := termdash.Run(ctx, tbx, cont, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}
}
