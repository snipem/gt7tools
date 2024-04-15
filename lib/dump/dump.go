package dump

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// WriteGT7Data writes compressed GT7Data slice to a file.
func WriteGT7Data(filename string, data []gt7.GTData) error {
	// Create a new file for writing.
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip writer.
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// Create a GOB encoder.
	encoder := gob.NewEncoder(gzipWriter)

	// Encode and write each GT7Data instance in the slice.
	for _, d := range data {
		if err := encoder.Encode(d); err != nil {
			return err
		}
	}

	return nil
}

// ReadGT7Data reads compressed GT7Data from a file.
func ReadGT7Data(filename string) ([]gt7.GTData, error) {
	var result []gt7.GTData

	// Read the whole file.
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Create a reader for the compressed data.
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Create a GOB decoder.
	decoder := gob.NewDecoder(reader)

	// Decode GT7Data until EOF.
	for {
		var d gt7.GTData
		if err := decoder.Decode(&d); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		result = append(result, d)
	}

	return result, nil
}

type GT7Dump struct {
	LastData          gt7.GTData
	data              []gt7.GTData
	gt7c              *gt7.GT7Communication
	DataSendFrequency time.Duration
}

func NewGT7Dump(filename string, gt7c *gt7.GT7Communication) (GT7Dump, error) {
	data, err := ReadGT7Data(filename)
	if err != nil {
		return GT7Dump{}, err
	}
	gt7d := GT7Dump{
		LastData:          gt7.GTData{},
		data:              data,
		gt7c:              gt7c,
		DataSendFrequency: 16 * time.Millisecond,
	}
	return gt7d, nil
}

func (gt7d *GT7Dump) Run() {

	for {
		for i := 0; i < len(gt7d.data); i++ {
			gt7d.LastData = gt7d.data[i]
			gt7d.gt7c.LastData = gt7d.LastData

			time.Sleep(gt7d.DataSendFrequency)
		}
		// Start over
	}

}

func NewRealOrDumpedGT7Connection(gt7c *gt7.GT7Communication, playstationIp string, dumpFilePath string) error {
	gt7c = gt7.NewGT7Communication(playstationIp)

	if dumpFilePath != "" {

		gt7dump, err := NewGT7Dump(dumpFilePath, gt7c)
		if err != nil {
			return fmt.Errorf("error loading dump file: %v", err)
		}
		log.Println("Using dump file: ", dumpFilePath)
		gt7dump.Run()

	} else {
		for {
			err := gt7c.Run()
			if err != nil {
				log.Printf("error running gt7c.Run(): %v", err)
			}
			log.Println("Sleeping 10 seconds before restarting gt7c.Run()")
			time.Sleep(10 * time.Second)
		}
	}
	return nil
}
