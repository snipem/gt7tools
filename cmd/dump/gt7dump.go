package main

import (
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/snipem/gt7tools/lib/dump"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Example usage:
	gt7c := gt7.NewGT7Communication("255.255.255.255")

	go gt7c.Run()

	// FIXME overwrite datafile initially

	filename := "datafile.gob.gz"

	deleteFileIfExists(filename)

	oldPackageId := int32(0)
	nrPackagesReceived := 0
	var data []gt7.GTData
	go func() {
		for {
			if gt7c.LastData.PackageID != oldPackageId {
				data = append(data, gt7c.LastData)
				nrPackagesReceived++
				fmt.Println("Received package:", nrPackagesReceived)
			}
			oldPackageId = gt7c.LastData.PackageID
			//time.Sleep(16 * time.Millisecond)
			// This will lead to missed packages
		}

	}()

	// Setup signal handling for Ctrl+C or abort signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for the signal
	fmt.Println("Waiting for Ctrl+C (SIGINT) or abort signal (SIGTERM) to write data to file...")
	<-sigCh

	// Write data to file.
	if err := dump.WriteGT7Data(filename, data); err != nil {
		fmt.Println("Error writing data:", err)
		return
	}
	fmt.Printf("\n%d packages written to file:%s\n", len(data), filename)

}

func deleteFileIfExists(filename string) {

	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			fmt.Println("Error deleting file:", err)
		}
	}

}
