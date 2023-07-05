package main

import (
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"time"
)

func main() {
	gt7c := gt7.NewGT7Communication("255.255.255.255")

	speech := htgotts.Speech{Folder: "audio", Language: voices.English}

	go gt7c.Run()
	dataBefore := gt7.GTData{}
	for true {
		if dataBefore.PackageID != gt7c.LastData.PackageID {
			checkForIdealDownshift(dataBefore, gt7c.LastData, speech)
		}
		dataBefore = gt7c.LastData
		time.Sleep(16 * time.Millisecond)
	}
}

func getCarName(id int32) string {
	return fmt.Sprintf("%d", id)
}

func log(dataNow gt7.GTData, suffix string) {
	fmt.Printf("L%d [%s] - %s\n", dataNow.CurrentLap, getCarName(dataNow.CarID), suffix)
}

func checkForIdealDownshift(dataBefore gt7.GTData, dataNow gt7.GTData, speech htgotts.Speech) {

	// Downshift
	if dataBefore.CurrentGear > dataNow.CurrentGear {
		revRatio := dataNow.RPM / float32(dataNow.RPMRevWarning)
		if revRatio >= 0.80 {
			print("\a") // bell
			log(dataNow, fmt.Sprintf("Poor Downshift %d > %d: %.2f, RPM Before: %.0f, RPM After: %.0f, REV Warning: %d", dataBefore.CurrentGear, dataNow.CurrentGear, revRatio, dataBefore.RPM, dataNow.RPM, dataNow.RPMRevWarning))
			//speech.Speak("Poor Downshift")
		}

	}

}
