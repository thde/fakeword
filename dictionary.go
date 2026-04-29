package fakeword

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

const defaultPrefixLength = MaxSequencesDefault + 1

// Dictionary stores words to be used to create a Generator.
type Dictionary struct {
	PrefixLength int
	counter      map[string]map[string]int
}

// Add words to a Dictionary.
func (w *Dictionary) Add(words ...string) *Dictionary {
	if w.PrefixLength == 0 {
		w.PrefixLength = defaultPrefixLength
	}

	for _, word := range words {
		word := strings.ToLower(strings.TrimSpace(word))
		word = prefix + word + suffix

		for i := 2; i <= w.PrefixLength; i++ {
			for _, substr := range splitToLength(word, i) {
				w.count(substr)
			}
		}
	}

	return w
}

// Read from an io.Reader and adds those words to a Dictionary.
// Lines prefixed with # are skipped.
func (w *Dictionary) Read(in io.Reader) *Dictionary {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		words := strings.Fields(line)
		w.Add(words...)
	}

	return w
}

// Generator returns a new Generator based on the words added
// to the dictionary.
func (w *Dictionary) Generator() Generator {
	probs := make(map[string]map[string]float32, len(w.counter))
	compiled := make(map[string][]outcome, len(w.counter))

	for ctx, counts := range w.counter {
		var sum int
		for _, c := range counts {
			sum += c
		}

		type pair struct {
			sym   byte
			count int
		}
		pairs := make([]pair, 0, len(counts))
		inner := make(map[string]float32, len(counts))
		for s, c := range counts {
			inner[s] = float32(c) / float32(sum)
			pairs = append(pairs, pair{s[0], c})
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].sym < pairs[j].sym })

		outcomes := make([]outcome, len(pairs))
		var cum float32
		for i, p := range pairs {
			cum += float32(p.count) / float32(sum)
			outcomes[i] = outcome{sym: p.sym, cum: cum}
		}
		if len(outcomes) > 0 {
			outcomes[len(outcomes)-1].cum = 1.0
		}

		probs[ctx] = inner
		compiled[ctx] = outcomes
	}

	return Generator{Probabilities: probs, compiled: compiled}
}

// count the amount of occurencies of a suffix.
func (w *Dictionary) count(substr string) {
	prefix := substr[:len(substr)-1]
	suffix := substr[len(substr)-1:]

	if w.counter == nil {
		w.counter = map[string]map[string]int{}
	}

	_, ok := w.counter[prefix]
	if !ok {
		w.counter[prefix] = map[string]int{}
	}

	w.counter[prefix][suffix]++
}

// splitToLength splits a string to substrings of length.
func splitToLength(s string, length int) []string {
	substrs := []string{}

	for i := 0; i <= len(s)-1; i++ {
		j := i + length
		if j > len(s) {
			continue
		}

		substrs = append(substrs, s[i:j])
	}

	return substrs
}
