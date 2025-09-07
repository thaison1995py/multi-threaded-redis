package data_structure

import (
	"testing"
)

func TestNewSimpleSet(t *testing.T) {
	key := "test_key"
	ss := NewSimpleSet(key)
	
	if ss.key != key {
		t.Errorf("Expected key %s, got %s", key, ss.key)
	}
	
	if ss.dict == nil {
		t.Error("Expected dict to be initialized, got nil")
	}
	
	if len(ss.dict) != 0 {
		t.Errorf("Expected empty dict, got %d items", len(ss.dict))
	}
}

func TestSimpleSet_Add(t *testing.T) {
	ss := NewSimpleSet("test")
	
	// Test adding single member
	added := ss.Add("member1")
	if added != 1 {
		t.Errorf("Expected 1 added member, got %d", added)
	}
	
	// Test adding duplicate member
	added = ss.Add("member1")
	if added != 0 {
		t.Errorf("Expected 0 added members (duplicate), got %d", added)
	}
	
	// Test adding multiple members
	added = ss.Add("member2", "member3", "member4")
	if added != 3 {
		t.Errorf("Expected 3 added members, got %d", added)
	}
	
	// Test adding mixed new and existing members
	added = ss.Add("member2", "member5", "member6")
	if added != 2 {
		t.Errorf("Expected 2 added members (member5, member6), got %d", added)
	}
	
	// Verify all members are present
	expectedMembers := []string{"member1", "member2", "member3", "member4", "member5", "member6"}
	members := ss.Members()
	if len(members) != len(expectedMembers) {
		t.Errorf("Expected %d members, got %d", len(expectedMembers), len(members))
	}
}

func TestSimpleSet_Rem(t *testing.T) {
	ss := NewSimpleSet("test")
	
	// Add some members first
	ss.Add("member1", "member2", "member3", "member4")
	
	// Test removing single member
	removed := ss.Rem("member1")
	if removed != 1 {
		t.Errorf("Expected 1 removed member, got %d", removed)
	}
	
	// Test removing non-existent member
	removed = ss.Rem("nonexistent")
	if removed != 0 {
		t.Errorf("Expected 0 removed members (non-existent), got %d", removed)
	}
	
	// Test removing multiple members
	removed = ss.Rem("member2", "member3", "member5")
	if removed != 2 {
		t.Errorf("Expected 2 removed members (member2, member3), got %d", removed)
	}
	
	// Verify remaining members
	members := ss.Members()
	if len(members) != 1 {
		t.Errorf("Expected 1 remaining member, got %d", len(members))
	}
	if members[0] != "member4" {
		t.Errorf("Expected remaining member to be 'member4', got %s", members[0])
	}
}

func TestSimpleSet_IsMember(t *testing.T) {
	ss := NewSimpleSet("test")
	
	// Test non-existent member
	if ss.IsMember("nonexistent") != 0 {
		t.Error("Expected non-existent member to return 0")
	}
	
	// Add members
	ss.Add("member1", "member2")
	
	// Test existing members
	if ss.IsMember("member1") != 1 {
		t.Error("Expected existing member1 to return 1")
	}
	if ss.IsMember("member2") != 1 {
		t.Error("Expected existing member2 to return 1")
	}
	
	// Test non-existent member after adding others
	if ss.IsMember("member3") != 0 {
		t.Error("Expected non-existent member3 to return 0")
	}
}

func TestSimpleSet_Members(t *testing.T) {
	ss := NewSimpleSet("test")
	
	// Test empty set
	members := ss.Members()
	if len(members) != 0 {
		t.Errorf("Expected empty members list, got %d items", len(members))
	}
	
	// Add members
	ss.Add("member1", "member2", "member3")
	
	// Test members list
	members = ss.Members()
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}
	
	// Verify all expected members are present
	expectedMembers := map[string]bool{
		"member1": true,
		"member2": true,
		"member3": true,
	}
	
	for _, member := range members {
		if !expectedMembers[member] {
			t.Errorf("Unexpected member: %s", member)
		}
	}
}

func TestSimpleSet_EdgeCases(t *testing.T) {
	ss := NewSimpleSet("test")
	
	// Test adding empty string
	added := ss.Add("")
	if added != 1 {
		t.Errorf("Expected 1 added empty string, got %d", added)
	}
	
	// Test removing empty string
	removed := ss.Rem("")
	if removed != 1 {
		t.Errorf("Expected 1 removed empty string, got %d", removed)
	}
	
	// Test adding same member multiple times
	ss.Add("duplicate", "duplicate", "duplicate")
	members := ss.Members()
	if len(members) != 1 {
		t.Errorf("Expected 1 unique member, got %d", len(members))
	}
	if members[0] != "duplicate" {
		t.Errorf("Expected 'duplicate', got %s", members[0])
	}
}
