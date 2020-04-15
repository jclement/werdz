package words

import (
	"strings"
	"testing"
)

const testJson = `
    {
      "word": {
        "value": "habanera",
        "definition": "slow and seductive Cuban dance"
      }
    }
    {
      "word": {
        "value": "habergeon",
        "definition": "sleeveless mail coat"
      }
    }
    {
      "word": {
        "value": "habilable",
        "definition": "capable of being clothed"
      }
	}
`

func TestBasics(t *testing.T) {
	rdr := strings.NewReader(testJson)

	var words WordSet
	words.Load(rdr)

	seenWords := make(map[string]bool)
	for i := 0; i < 100; i++ {
		w := words.Random().Word
		if _, ok := seenWords[w]; !ok {
			seenWords[w] = true
		}
		if len(seenWords) == 3 {
			return
		}
	}

	t.Error("Not working!")
}
