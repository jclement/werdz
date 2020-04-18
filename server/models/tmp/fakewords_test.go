package fakewords

import (
	"strings"
	"testing"
)

const testJson = `
    {
      "word": "Liminary"
    }
    {
      "word": "Molish"
    }
    {
      "word": "Scrutable"
    }
`

func TestBasics(t *testing.T) {
	rdr := strings.NewReader(testJson)

	var words FakeWordSet
	words.Load(rdr)

	seenWords := make(map[string]bool)
	for i := 0; i < 100; i++ {
		w := words.Random()
		if _, ok := seenWords[w]; !ok {
			seenWords[w] = true
		}
		if len(seenWords) == 3 {
			return
		}
	}

	t.Error("Not working!")
}
