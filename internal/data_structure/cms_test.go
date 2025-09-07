package data_structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCMS(t *testing.T) {
	w := uint32(100)
	d := uint32(5)
	cms := CreateCMS(w, d)

	assert.NotNil(t, cms)
	assert.Equal(t, w, cms.width)
	assert.Equal(t, d, cms.depth)
	assert.Len(t, cms.counter, int(d))
	assert.Len(t, cms.counter[0], int(w))
	assert.Equal(t, uint32(0), cms.counter[0][0])
}

func TestCalcCMSDim(t *testing.T) {
	errRate := 0.01
	errProb := 0.001

	w, d := CalcCMSDim(errRate, errProb)

	// Expected values based on the formulas:
	// w = ceil(2.0 / errRate) = ceil(2.0 / 0.01) = 200
	// d = ceil(log10(errProb) / Log10PointFive) = ceil(log10(0.001) / -0.30102999566) = ceil(-3 / -0.30102999566) = ceil(9.96) = 10
	assert.Equal(t, uint32(200), w)
	assert.Equal(t, uint32(10), d)
}

func TestCMS_IncrByAndCount(t *testing.T) {
	w := uint32(100)
	d := uint32(5)
	cms := CreateCMS(w, d)

	item1 := "test_item_1"
	item2 := "test_item_2"

	// Test IncrBy
	cms.IncrBy(item1, 5)
	assert.True(t, cms.Count(item1) >= 5)

	cms.IncrBy(item1, 10)
	assert.True(t, cms.Count(item1) >= 15)

	cms.IncrBy(item2, 3)
	assert.True(t, cms.Count(item2) >= 3)

	// Test Count for non-existent item
	assert.Equal(t, uint32(0), cms.Count("non_existent_item"))

	// Test with multiple increments and check for minimum count property
	for i := 0; i < 1000; i++ {
		cms.IncrBy("another_item", 1)
	}
	assert.True(t, cms.Count("another_item") >= 1000)
}
