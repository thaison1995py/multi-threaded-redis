package data_structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBPlusTree(t *testing.T) {
	degree := 3
	bt := NewBPlusTree(degree)

	assert.NotNil(t, bt)
	assert.Equal(t, degree, bt.Degree)
	assert.NotNil(t, bt.Root)
	assert.True(t, bt.Root.IsLeaf)
	assert.Len(t, bt.Root.Items, 0)
	assert.Len(t, bt.Root.Children, 0)
}

func TestBPlusTree_AddAndScore(t *testing.T) {
	degree := 3
	bt := NewBPlusTree(degree)

	bt.Add(10.0, "member1")
	bt.Add(20.0, "member2")
	bt.Add(5.0, "member3")

	score, found := bt.Score("member1")
	assert.True(t, found)
	assert.Equal(t, 10.0, score)

	score, found = bt.Score("member2")
	assert.True(t, found)
	assert.Equal(t, 20.0, score)

	score, found = bt.Score("member3")
	assert.True(t, found)
	assert.Equal(t, 5.0, score)

	_, found = bt.Score("non_existent_member")
	assert.False(t, found)
}

func TestBPlusTree_GetRank(t *testing.T) {
	degree := 3
	bt := NewBPlusTree(degree)

	bt.Add(10.0, "memberA")
	bt.Add(20.0, "memberB")
	bt.Add(5.0, "memberC")
	bt.Add(15.0, "memberD")

	// The ranks depend on the order in which elements are stored internally.
	// For a simple B+ tree without explicit rank tracking, we'll assume a sorted order for verification.
	// memberC (5.0), memberA (10.0), memberD (15.0), memberB (20.0)
	assert.Equal(t, 0, bt.GetRank("memberC"))
	assert.Equal(t, 1, bt.GetRank("memberA"))
	assert.Equal(t, 2, bt.GetRank("memberD"))
	assert.Equal(t, 3, bt.GetRank("memberB"))
	assert.Equal(t, -1, bt.GetRank("non_existent_member"))
}
