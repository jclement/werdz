package player

import "testing"

func TestTestFunc(t *testing.T) {
	if Test() != "hello" {
		t.Errorf("Nope!")
	}
}
