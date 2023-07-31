package lib

import (
	"fmt"
	"github.com/0xcafed00d/joystick"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"time"
)

type FakeGT7Communication struct {
	LastData gt7.GTData
}

func NewFakeGT7Communication() *FakeGT7Communication {

	return &FakeGT7Communication{
		LastData: gt7.GTData{},
	}

}

func (gt7c *FakeGT7Communication) Run() {

	js, err := joystick.Open(0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Joystick Name: %s", js.Name())
	fmt.Printf("   Axis Count: %d", js.AxisCount())
	fmt.Printf(" Button Count: %d", js.ButtonCount())

	const maxValue = 32768

	for true {
		state, err := js.Read()
		if err != nil {
			panic(err)
		}

		gt7c.LastData = gt7.GTData{
			Throttle: float32(state.AxisData[js.AxisCount()-1]) / float32(maxValue) * 100,
			Brake:    float32(state.AxisData[js.AxisCount()-2]) / float32(maxValue) * 100,
		}
		time.Sleep(100 * time.Millisecond)
	}

	js.Close()
}
