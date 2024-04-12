package dump

import (
	"fmt"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGT7Dump_Run(t *testing.T) {
	gt7c, err := NewGT7Dump("../../datafile.gob.gz", nil)
	assert.NoError(t, err)
	go gt7c.Run()

	lastPackageId := int32(0)
	packagesSeen := 0
	for packagesSeen < 100 {
		if lastPackageId != gt7c.LastData.PackageID {
			fmt.Printf("PackageID: %d, CurrentLap: %d, RPM: %0.f\n", gt7c.LastData.PackageID, gt7c.LastData.CurrentLap, gt7c.LastData.RPM)
			packagesSeen++
			lastPackageId = gt7c.LastData.PackageID
		}
	}
}

func TestReadGT7Data(t *testing.T) {

	data, err := ReadGT7Data("../../datafile.gob.gz")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}

func TestWriteGT7Data(t *testing.T) {

	err := WriteGT7Data("test_datafile.gob.gz", []gt7.GTData{})
	if err != nil {
		t.Error(err)
	}
}
