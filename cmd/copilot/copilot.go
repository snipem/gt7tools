package main

import (
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"log"
	"time"
)

func speek(s string) {
	speech := htgotts.Speech{Folder: "audio", Language: voices.English}
	err := speech.Speak(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func main() {
	gt7c := gt7.NewGT7Communication("255.255.255.255")
	go gt7c.Run()

	dataBefore := gt7c.LastData
	warnedInGear := uint8(0)
	for true {

		//if dataBefore.CurrentPosition != gt7c.LastData.CurrentPosition {
		//	speek(fmt.Sprintf("P%d", gt7c.LastData.CurrentLap))
		//}

		if gt7c.LastData.CurrentLap != -1 && dataBefore.CurrentLap != gt7c.LastData.CurrentLap && gt7c.LastData.TotalLaps != 0 {
			if gt7c.LastData.TotalLaps == gt7c.LastData.CurrentLap {
				speek("Final Lap!")
			}
			if gt7c.LastData.TotalLaps < gt7c.LastData.CurrentLap {
				speek("Finish!")
			} else {
				speek(fmt.Sprintf("Lap %d of %d", gt7c.LastData.CurrentLap, gt7c.LastData.TotalLaps))
			}
		}

		maxGear := getMaxGear(gt7c.LastData)

		if gt7c.LastData.CurrentGear != warnedInGear && gt7c.LastData.Throttle == 100 &&
			gt7c.LastData.CurrentGear != maxGear &&
			dataBefore.Boost > gt7c.LastData.Boost &&
			dataBefore.CurrentGear == gt7c.LastData.CurrentGear {
			fmt.Printf("Current Gear: %d, RPM to Shift: %d\n", gt7c.LastData.CurrentGear, int(gt7c.LastData.RPM))
			fmt.Printf("Boost Before: %.5f, Boost Now: %.5f\n", dataBefore.Boost, gt7c.LastData.Boost)
			speek("Shift!")
			warnedInGear = gt7c.LastData.CurrentGear
		}

		dataBefore = gt7c.LastData

		time.Sleep(16 * time.Millisecond)
	}
}

func getMaxGear(data gt7.GTData) uint8 {

	if data.Gear1 == 0 {
		return 0
	} else if data.Gear2 == 0 {
		return 1
	} else if data.Gear3 == 0 {
		return 2
	} else if data.Gear4 == 0 {
		return 3
	} else if data.Gear5 == 0 {
		return 4
	} else if data.Gear6 == 0 {
		return 5
	} else if data.Gear7 == 0 {
		return 6
	} else if data.Gear8 == 0 {
		return 7
	}

	return 0

}
