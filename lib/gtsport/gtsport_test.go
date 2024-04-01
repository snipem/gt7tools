package gtsport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOnlineResult(t *testing.T) {
	result, err := GetOnlineResult(9295, 0)
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestGetFolder(t *testing.T) {
	folder, err := GetFolder(6, 249)
	assert.NotNil(t, folder)
	assert.NoError(t, err)
	assert.NotNil(t, folder.Result[0].Parameters.Event.EventDescriptionKey)
}

func TestGetActiveTimeTrials(t *testing.T) {
	timetrials, err := GetActiveTimeTrials()
	//assert.Printf("%v", timetrials)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(timetrials), 1)
}
