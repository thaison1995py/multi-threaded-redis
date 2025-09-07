package data_structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSortedSet(t *testing.T) {
	degree := 3
	ss := NewSortedSet(degree)

	assert.NotNil(t, ss)
	assert.NotNil(t, ss.Tree)
	assert.Equal(t, degree, ss.Tree.Degree)
	assert.NotNil(t, ss.MemberScores)
	assert.Len(t, ss.MemberScores, 0)
}

func TestSortedSet_Add(t *testing.T) {
	degree := 3
	ss := NewSortedSet(degree)

	// Add new members
	added := ss.Add(10.0, "memberA")
	assert.Equal(t, 1, added)
	assert.Equal(t, 10.0, ss.MemberScores["memberA"])

	added = ss.Add(20.0, "memberB")
	assert.Equal(t, 1, added)
	assert.Equal(t, 20.0, ss.MemberScores["memberB"])

	// Update an existing member's score
	added = ss.Add(15.0, "memberA")
	assert.Equal(t, 1, added)
	assert.Equal(t, 15.0, ss.MemberScores["memberA"])

	// Check tree size (indirectly)
	score, found := ss.Tree.Score("memberA")
	assert.True(t, found)
	assert.Equal(t, 15.0, score)
}

func TestSortedSet_GetScore(t *testing.T) {
	degree := 3
	ss := NewSortedSet(degree)

	ss.Add(10.0, "member1")
	ss.Add(20.0, "member2")

	score, found := ss.GetScore("member1")
	assert.True(t, found)
	assert.Equal(t, 10.0, score)

	score, found = ss.GetScore("member2")
	assert.True(t, found)
	assert.Equal(t, 20.0, score)

	_, found = ss.GetScore("non_existent_member")
	assert.False(t, found)
}

func TestSortedSet_GetRank(t *testing.T) {
	degree := 3
	ss := NewSortedSet(degree)

	ss.Add(10.0, "memberA")
	ss.Add(20.0, "memberB")
	ss.Add(5.0, "memberC")
	ss.Add(15.0, "memberD")

	// The ranks depend on the order in which elements are stored internally (lowest score gets rank 0).
	// memberC (5.0), memberA (10.0), memberD (15.0), memberB (20.0)
	assert.Equal(t, 0, ss.GetRank("memberC"))
	assert.Equal(t, 1, ss.GetRank("memberA"))
	assert.Equal(t, 2, ss.GetRank("memberD"))
	assert.Equal(t, 3, ss.GetRank("memberB"))
	assert.Equal(t, -1, ss.GetRank("non_existent_member"))
}
