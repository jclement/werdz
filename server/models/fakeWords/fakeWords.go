package fakewords

import (
	"encoding/json"
	"io"
	"math/rand"
)

// FakeWordSet represents a set of words and definitions
type FakeWordSet struct {
	Words []string
}

type sourceRec struct {
	Word string `json:"word"`
}

// Load loads a JSON stream full of words into the FakeWordSet
func (w *FakeWordSet) Load(reader io.Reader) {
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
		w.Words = append(w.Words, rec.Word)
	}
}

// Random returns a random word from the wordset
func (w *FakeWordSet) Random() string {
	return w.Words[rand.Intn(len(w.Words))]
}
