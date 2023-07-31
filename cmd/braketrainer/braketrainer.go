package main

import (
	"encoding/csv"
	"fmt"
	"github.com/snipem/gt7-utils/lib"
	"log"
	"os"
	"time"
)

func main() {
	gt7c := lib.NewFakeGT7Communication()
	go gt7c.Run()

	gameRuns := true

	for gameRuns {
		start := time.Now()

		checkBrakes(&gt7c.LastData.Brake, 5)
		checkBrakes(&gt7c.LastData.Brake, 50)
		checkBrakes(&gt7c.LastData.Brake, 99)

		i := 100

		for i != 0 {
			checkBrakes(&gt7c.LastData.Brake, float32(i))
			i--
		}
		end := time.Now()
		fmt.Printf("Done! Time: %s\n", end.Sub(start))
		score := end.Sub(start).Seconds()

		fmt.Printf("Get ready for next round ...")
		writeScore(score)

		time.Sleep(5000 * time.Millisecond)

	}
}

func writeScore(score float64) {

	csvfile, err := os.OpenFile("braketrainer.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvfile)
	timestring := fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
	scorestring := fmt.Sprintf("%d", int(score))
	err = csvwriter.Write([]string{timestring, scorestring})
	if err != nil {
		log.Fatal(err)
	}

	csvwriter.Flush()

	csvfile.Close()

}

func checkBrakes(brake *float32, target float32) {
	brakeTargetNotMet := true

	for brakeTargetNotMet {

		err := lib.Flush()
		if err != nil {
			log.Fatal(err)
		}

		intBrake := int(*brake)
		intTarget := int(target)
		fmt.Printf("\rCurrent: %d Target: %d\n", intBrake, intTarget)

		if intBrake == intTarget {
			brakeTargetNotMet = false
		}
		time.Sleep(10 * time.Millisecond)
	}

}
