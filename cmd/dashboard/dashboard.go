package main

import (
	"context"
	"fmt"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/snipem/gt7-utils/lib"
	"math"
	"time"
)

// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary linechartdemo displays a linechart widget.
// Exist when 'q' is pressed.
var gt7c = &gt7.GT7Communication{}

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
			i = (i + 1) % len(inputs)
			if err := lc.Series("throttle", convertIntSliceToFloatSlice(takeLastN(history.Throttle, 1000)),
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(64))),
			); err != nil {
				panic(err)
			}

			if err := lc.Series("braking", convertIntSliceToFloatSlice(takeLastN(history.Brake, 1000)),
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(160))),
			); err != nil {
				panic(err)
			}

			//i2 := (i + 100) % len(inputs)
			//rotated2 := append(inputs[i2:], inputs[:i2]...)
			//if err := lc.Series("second", rotated2, linechart.SeriesCellOpts(cell.FgColor(cell.ColorWhite))); err != nil {
			//	panic(err)
			//}

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

	history := &lib.History{}

	go lib.UpdateHistory(gt7c, history)

	const redrawInterval = 500 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorBlack)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
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

	flText, _ := text.New()
	frText, _ := text.New()
	rlText, _ := text.New()
	rrText, _ := text.New()

	go playLineChart(ctx, lc, history, redrawInterval/3)
	go playBarChart(ctx, bc, redrawInterval/3)
	go updateTires(ctx, flText, frText, rlText, rrText, redrawInterval/3)

	tbx, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer tbx.Close()

	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(
			97,
			grid.ColWidthPerc(
				30,
				grid.RowHeightPerc(50, grid.Widget(flText, container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalBottom))),
				grid.RowHeightPerc(50, grid.Widget(rlText, container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalBottom))),
			),
			grid.ColWidthPerc(30, grid.Widget(bc)),
			grid.ColWidthPerc(
				30,
				grid.RowHeightPerc(50, grid.Widget(frText, container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalBottom))),
				grid.RowHeightPerc(50, grid.Widget(rrText, container.AlignHorizontal(align.HorizontalCenter), container.AlignVertical(align.VerticalBottom))),
			),
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

func updateTires(ctx context.Context, flText *text.Text, frText *text.Text, rlText *text.Text, rrText *text.Text, delay time.Duration) {

	const max = 100

	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			flText.Reset()
			frText.Reset()
			rlText.Reset()
			rrText.Reset()

			flText.Write(getTireString(gt7c.LastData.TyreTempFL), text.WriteCellOpts(getColorOptForTemp(gt7c.LastData.TyreTempFL)))
			frText.Write(getTireString(gt7c.LastData.TyreTempFR), text.WriteCellOpts(getColorOptForTemp(gt7c.LastData.TyreTempFR)))
			rlText.Write(getTireString(gt7c.LastData.TyreTempRL), text.WriteCellOpts(getColorOptForTemp(gt7c.LastData.TyreTempRL)))
			rrText.Write(getTireString(gt7c.LastData.TyreTempRR), text.WriteCellOpts(getColorOptForTemp(gt7c.LastData.TyreTempRR)))

		case <-ctx.Done():
			return
		}
	}
}

func getTireString(temp float32) string {
	return fmt.Sprintf("  %2.0f°  \n       \n       \n       \n       ", temp)
}

func getColorOptForTemp(temp float32) cell.Option {

	if temp >= 78 {
		return cell.BgColor(cell.ColorRed)
	} else if temp >= 72 {
		return cell.BgColor(cell.ColorGreen)
	}

	return cell.BgColor(cell.ColorBlue)

}
