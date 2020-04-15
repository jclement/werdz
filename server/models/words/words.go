package words

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
)

// Word represets a word/definition pair
type Word struct {
	Word       string
	Definition string
}

// WordSet represents a set of words and definitions
type WordSet struct {
	Words []Word
}

type sourceWord struct {
	Value      string `json:"value"`
	Definition string `json:"definition"`
}

type sourceRec struct {
	Word sourceWord `json:"word"`
}

// Load loads a JSON stream full of words into the WordSet
func (w *WordSet) Load(reader io.Reader) {
	var rec sourceRec
	decoder := json.NewDecoder(reader)
	for {
		err := decoder.Decode(&rec)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		w.Words = append(w.Words, Word{
			Word:       rec.Word.Value,
			Definition: rec.Word.Definition,
		})
	}
}

// Random returns a random word from the wordset
func (w *WordSet) Random() Word {
	return w.Words[rand.Intn(len(w.Words))]
}
