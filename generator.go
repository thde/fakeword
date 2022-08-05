// Package fakeword allows to generate fake words.
// Adding words of a certain language, allows to
// generate language like words.
package fakeword

import (
	"math/rand"
	"strings"
	"time"
)

const (
	MaxSequencesDefault = 2

	prefix = "^"
	suffix = "$"
)

type (
	// Generator allows to generate fake words.
	Generator struct {
		// Probabilities stores the probabilities of characters following on a string.
		Probabilities map[string]map[string]float32
		// MaxSequences defines how far back the algorithm looks
		// to predict the next character. A smaller value icreases randomness
		// and a higher value creates words that are closer to the dictionary words.
		// The default value is defined in MaxSequencesDefault.
		MaxSequences int
	}
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// Word generates a fake word with arbitrary length.
func (g Generator) Word() string {
	if len(g.Probabilities) == 0 {
		return ""
	}
	if g.MaxSequences == 0 {
		g.MaxSequences = MaxSequencesDefault
	}

	character := prefix
	word := ""
	characters := []string{}

	for character != suffix {
		characters = append(characters, character)
		if len(characters) > g.MaxSequences {
			characters = characters[1:]
		}
		var nextAccumedProbs map[string]float32
		n := 0
		for {
			str := strings.Join(characters[n:], "")
			nextAccumedProbs = g.Probabilities[str]
			n += 1
			if nextAccumedProbs != nil || n >= len(characters) {
				break
			}
		}
		nextCharacter := ""
		r := random.Float32()
		probability := float32(0)
		for ch, prob := range nextAccumedProbs {
			nextCharacterCandidate := ch
			probability += prob
			if r <= probability {
				nextCharacter = nextCharacterCandidate
				break
			}
		}
		if nextCharacter != suffix {
			word += nextCharacter
		}
		character = nextCharacter
	}
	return word
}

// WordWithDistance returns a fake word with variable length.
// Running this function can be quite costly, as it will just
// generate fake words until one with the correct length appears.
func (g Generator) WordWithDistance(min int, max int) string {
	fakeWord := ""
	for min > len(fakeWord) || len(fakeWord) > max {
		fakeWord = g.Word()
	}
	return fakeWord
}
