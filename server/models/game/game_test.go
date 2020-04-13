package game

import (
	"strings"
	"testing"
)

func TestGameIDUnique(t *testing.T) {
	seen := make(map[GameID]bool)
	for i := 0; i < 100; i++ {
		id := generateGameID()
		if _, exists := seen[id]; exists {
			t.Error("Generated Game IDs must be unique")
			return
		}
		seen[id] = true
	}
}

func TestGameIDSize(t *testing.T) {
	if len(generateGameID()) != 5 {
		t.Error("Generated Game IDs must be 5 characters")
	}
}

func TestGameIDCase(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := string(generateGameID())
		if id != strings.ToUpper(id) {
			t.Error("Generated Game ID must be upper case")
			return
		}
	}
}

func TestPlayerIDUnique(t *testing.T) {
	seen := make(map[PlayerID]bool)
	for i := 0; i < 100; i++ {
		id := GeneratePlayerID()
		if _, exists := seen[id]; exists {
			t.Error("Generated Player IDs must be unique")
			return
		}
		seen[id] = true
	}
}

func TestPlayerIDSize(t *testing.T) {
	if len(GeneratePlayerID()) != 20 {
		t.Error("Generated Player IDs must be 20 characters")
	}
}
