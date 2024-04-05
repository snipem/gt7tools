package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
