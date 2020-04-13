package game

import (
	"strings"
	"testing"
)

func TestRoomIDUnique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateRoomID()
		if _, exists := seen[id]; exists {
			t.Error("Generated Room IDs must be unique")
			return
		}
		seen[id] = true
	}
}

func TestRoomIDSize(t *testing.T) {
	if len(generateRoomID()) != 5 {
		t.Error("Generated Room IDs must be 5 characters")
	}
}

func TestRoomIDCase(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := generateRoomID()
		if id != strings.ToUpper(id) {
			t.Error("Generated Room ID must be upper case")
			return
		}
	}
}
