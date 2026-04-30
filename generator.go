// Package fakeword allows to generate fake words.
// Adding words of a certain language, allows to
// generate language like words.
//
// The generator is character-based and assumes ASCII input.
// Multibyte runes will be processed byte-wise and produce
// nonsensical contexts.
package fakeword // import "thde.io/fakeword"

import (
	"math/rand/v2"
	"sort"
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

		// compiled is the cumulative-probability form of Probabilities,
		// kept for fast deterministic sampling. Populated by
		// Dictionary.Generator. When nil, Word builds outcomes on demand
		// from Probabilities (slower but still correct).
		compiled map[string][]outcome
	}

	// outcome is a possible next-symbol with the cumulative probability
	// up to and including this outcome within its context. Outcomes for
	// a context are sorted by sym ascending so cum is monotonic.
	outcome struct {
		sym byte
		cum float32
	}
)

// percentage converts the uint32 to a float32 in the half-open interval [0.0,1.0).
// https://cs.opensource.google/go/go/+/refs/tags/go1.22.0:src/math/rand/v2/rand.go;l=211
func percentage(n uint32) float32 {
	// There are exactly 1<<24 float32s in [0,1). Use Intn(1<<24) / (1<<24).
	return float32(n<<8>>8) / (1 << 24)
}

// next samples the next symbol from the model given the current
// buffer. If filter is non-nil it transforms the outcome set
// before sampling (used to suppress the suffix marker while
// WordWithDistance is below its minimum length).
// Returns 0, false if no outcome was available.
func (g Generator) next(buf []byte, filter func([]outcome) []outcome) (byte, bool) {
	maxSeq := g.MaxSequences
	if maxSeq == 0 {
		maxSeq = MaxSequencesDefault
	}
	start := max(len(buf)-maxSeq, 0)
	outcomes := g.contextOutcomes(buf[start:])
	if filter != nil {
		outcomes = filter(outcomes)
	}
	if len(outcomes) == 0 {
		return 0, false
	}
	randomFunc := g.Random
	if randomFunc == nil {
		randomFunc = rand.Uint32
	}
	r := percentage(randomFunc())
	idx := sort.Search(len(outcomes), func(i int) bool {
		return outcomes[i].cum >= r
	})
	if idx == len(outcomes) {
		idx = len(outcomes) - 1
	}
	return outcomes[idx].sym, true
}

// Word generates a fake word with arbitrary length.
func (g Generator) Word() string {
	if len(g.Probabilities) == 0 && len(g.compiled) == 0 {
		return ""
	}
	buf := []byte{prefix[0]}
	for {
		sym, ok := g.next(buf, nil)
		if !ok || sym == suffix[0] {
			break
		}
		buf = append(buf, sym)
	}
	return string(buf[1:])
}

// contextOutcomes returns the outcome distribution for the longest
// suffix of window that has a known context, applying stupid-backoff.
func (g Generator) contextOutcomes(window []byte) []outcome {
	for s := range window {
		key := string(window[s:])
		if g.compiled != nil {
			if oc, ok := g.compiled[key]; ok {
				return oc
			}
			continue
		}
		if probs, ok := g.Probabilities[key]; ok {
			return compileContext(probs)
		}
	}
	return nil
}

// compileContext converts a per-context probability map into a sorted
// cumulative-probability slice for binary-search sampling.
func compileContext(probs map[string]float32) []outcome {
	type pair struct {
		sym  byte
		prob float32
	}
	pairs := make([]pair, 0, len(probs))
	for s, p := range probs {
		if len(s) == 0 {
			continue
		}
		pairs = append(pairs, pair{s[0], p})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].sym < pairs[j].sym })

	outcomes := make([]outcome, len(pairs))
	var cum float32
	for i, p := range pairs {
		cum += p.prob
		outcomes[i] = outcome{sym: p.sym, cum: cum}
	}
	if len(outcomes) > 0 {
		outcomes[len(outcomes)-1].cum = 1.0
	}
	return outcomes
}

// WordWithDistance returns a fake word whose length is in [min, max].
// It conditions termination on length: the suffix marker is suppressed
// while the word is shorter than min, and the loop hard-stops at max.
//
// If a context has no non-suffix outcome before min is reached the
// word ends early; if no context terminates naturally before max the
// word is truncated.
func (g Generator) WordWithDistance(minLen, maxLen int) string {
	if minLen < 0 {
		minLen = 0
	}
	if max < min {
		max = min
	}
	buf := []byte{prefix[0]}
	for len(buf)-1 < max {
		var filter func([]outcome) []outcome
		if len(buf)-1 < min {
			filter = withoutSuffix
		}
		sym, ok := g.next(buf, filter)
		if !ok || sym == suffix[0] {
			break
		}
		buf = append(buf, sym)
	}
	return string(buf[1:])
}

// withoutSuffix returns the input outcomes with the suffix marker
// removed and the remaining cumulative probabilities renormalized.
// Returns nil if the input contained only the suffix marker.
func withoutSuffix(outcomes []outcome) []outcome {
	sIdx := -1
	for i, oc := range outcomes {
		if oc.sym == suffix[0] {
			sIdx = i
			break
		}
	}
	if sIdx == -1 {
		return outcomes
	}
	if len(outcomes) == 1 {
		return nil
	}

	var sMass float32
	if sIdx == 0 {
		sMass = outcomes[0].cum
	} else {
		sMass = outcomes[sIdx].cum - outcomes[sIdx-1].cum
	}
	keptMass := 1.0 - sMass
	if keptMass <= 0 {
		return nil
	}

	filtered := make([]outcome, 0, len(outcomes)-1)
	var prev, cum float32
	for i, oc := range outcomes {
		diff := oc.cum - prev
		prev = oc.cum
		if i == sIdx {
			continue
		}
		cum += diff / keptMass
		filtered = append(filtered, outcome{sym: oc.sym, cum: cum})
	}
	filtered[len(filtered)-1].cum = 1.0
	return filtered
}
