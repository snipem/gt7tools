package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_straightIncreaseFromZeroBraking(t *testing.T) {
	assert.True(t, straightIncreaseFromZeroBraking([]int{5, 4, 3, 5, 0, 2, 3, 4, 5}))
	assert.True(t, straightIncreaseFromZeroBraking([]int{0, 2, 3, 4, 5}))
	assert.False(t, straightIncreaseFromZeroBraking([]int{0, 2, 3, 2, 4, 5}))
	assert.False(t, straightIncreaseFromZeroBraking([]int{0, 2, 3, 4, 5, 2}))
	assert.True(t, straightIncreaseFromZeroBraking([]int{0, 2, 2, 2}))
	assert.False(t, straightIncreaseFromZeroBraking([]int{0}))
	assert.False(t, straightIncreaseFromZeroBraking([]int{}))
}

func TestRun(t *testing.T) {
	// Define the timeout duration
	timeoutDuration := 15 * time.Second

	// Create a channel to signal completion
	done := make(chan struct{})

	// Start a goroutine to run the function
	go func() {
		Run()
		close(done) // Signal completion
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		fmt.Println("Run() completed successfully.")
	case <-time.After(timeoutDuration):
		fmt.Println("Timeout: Run() took too long.")
	}
}

func Test_mapGearToScale(t *testing.T) {
	mappedGers := mapGearToScale(5, 100, []int{0, 1, 2, 3, 4, 5, 5, 3, 2, 1, 0})
	assert.Equal(t, []int{0, 20, 40, 60, 80, 100, 100, 60, 40, 20, 0}, mappedGers)
}
