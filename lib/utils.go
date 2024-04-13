package lib

import (
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"time"
)

type History struct {
	Throttle   []int
	Brake      []int
	Gear       []int
	Position_Z []float32
}

// UpdateHistory updates the history in an endless loop with throttle and breaking information
// it waits for 16 milliseconds which is equivalent to the framerate (1000ms / 60fps)
func UpdateHistory(gt7c *gt7.GT7Communication, history *History) {
	oldPackageId := int32(0)

	for {
		if gt7c.LastData.PackageID != oldPackageId {

			throttle := int(gt7c.LastData.Throttle)
			brake := int(gt7c.LastData.Brake)
			gear := int(gt7c.LastData.CurrentGear)
			positionZ := gt7c.LastData.PositionZ
			history.Throttle = append(history.Throttle, throttle)
			history.Brake = append(history.Brake, brake)
			history.Gear = append(history.Gear, gear)
			history.Position_Z = append(history.Position_Z, positionZ)
			//fmt.Printf("package id: %d, brake: %d, throttle: %d\n", gt7c.LastData.PackageID, brake, throttle)
			//fmt.Printf("Got %d packets\n", len(history.Throttle))
		}
		oldPackageId = gt7c.LastData.PackageID
		time.Sleep(16 * time.Millisecond)
	}
}

func Flush() error {
	_, err := fmt.Print("\033[H\033[2J")
	return err
}
