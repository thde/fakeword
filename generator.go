// Package fakeword allows to generate fake words.
// Adding words of a certain language, allows to
// generate language like words.
package fakeword // import "thde.io/fakeword"

import (
	"math/rand/v2"
	"strings"
)

const (
	// MaxSequencesDefault contains the default for Generator.MaxSequences.
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

		// Random should return a 32-bit value as a uint32.
		// Uses math/rand/v2's Uint32 if Random is nil.
		Random func() uint32
	}
)

// percentage converts the uint32 to a float32 in the half-open interval [0.0,1.0).
// https://cs.opensource.google/go/go/+/refs/tags/go1.22.0:src/math/rand/v2/rand.go;l=211
func percentage(n uint32) float32 {
	// There are exactly 1<<24 float32s in [0,1). Use Intn(1<<24) / (1<<24).
	return float32(n<<8>>8) / (1 << 24)
}

// Word generates a fake word with arbitrary length.
func (g Generator) Word() string {
	if len(g.Probabilities) == 0 {
		return ""
	}
	if g.MaxSequences == 0 {
		g.MaxSequences = MaxSequencesDefault
	}

	randomFunc := g.Random
	if randomFunc == nil {
		randomFunc = rand.Uint32
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
			n++
			if nextAccumedProbs != nil || n >= len(characters) {
				break
			}
		}

		nextCharacter := ""
		target := percentage(randomFunc())
		probability := float32(0)
		for ch, prob := range nextAccumedProbs {
			nextCharacterCandidate := ch
			probability += prob
			if target <= probability {
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
