package lib

import (
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"time"
)

type History struct {
	Throttle []int
	Brake    []int
}

// UpdateHistory updates the history in an endless loop with throttle and breaking information
// it waits for 16 milliseconds which is equivalent to the framerate (1000ms / 60fps)
func UpdateHistory(gt7c *gt7.GT7Communication, history *History) {
	oldPackageId := int32(0)

	for true {
		if gt7c.LastData.PackageID != oldPackageId {

			throttle := int(gt7c.LastData.Throttle)
			brake := int(gt7c.LastData.Brake)
			history.Throttle = append(history.Throttle, throttle)
			history.Brake = append(history.Brake, brake)
			//fmt.Printf("package id: %d, brake: %d, throttle: %d\n", gt7c.LastData.PackageID, brake, throttle)
			//fmt.Printf("Got %d packets\n", len(history.Throttle))
		}
		oldPackageId = gt7c.LastData.PackageID
		time.Sleep(16 * time.Millisecond)
	}
}
