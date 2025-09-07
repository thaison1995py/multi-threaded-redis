package data_structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBloom_Exist(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	b.Add("a")
	b.Add("b")
	assert.EqualValues(t, 10, b.Entries)
	assert.EqualValues(t, 0.01, b.Error)
	assert.True(t, b.Exist("a"))
	assert.True(t, b.Exist("b"))
	assert.False(t, b.Exist("c"))
	assert.False(t, b.Exist("d"))
}

func TestBloom_CalcHash(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	x := b.CalcHash("abcdef")
	y := b.CalcHash("abcdef")
	assert.EqualValues(t, x.a, y.a)
	assert.EqualValues(t, x.b, y.b)
}

func TestBloom_AddHash(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	hash := b.CalcHash("abcdef")
	b.AddHash(hash)
	assert.True(t, b.ExistHash(hash))
	assert.True(t, b.Exist("abcdef"))
}

func TestBloom_CreateBloomFilter(t *testing.T) {
	b := CreateBloomFilter(100, 0.01)
	assert.EqualValues(t, 100, b.Entries)
	assert.EqualValues(t, 0.01, b.Error)
	// Expected bits and bytes are calculated based on the formulas in bloom.go
	// bitsFloat = (entries * bitPerEntry) = 100 * 9.58496 = 958.496
	// bloom.bits = ceil(bitsFloat / 64.0) * 64 = ceil(958.496 / 64.0) * 64 = 15 * 64 = 960
	// bloom.bytes = bits / 8 = 960 / 8 = 120
	assert.EqualValues(t, 960, b.bits)
	assert.EqualValues(t, 120, b.bytes)
	assert.EqualValues(t, 7, b.Hashes) // hashes = bitPerEntry * ln(2) = 9.58496 * 0.693 = 6.639 -> ceil(6.639) = 7
}
