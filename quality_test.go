package fakeword

import (
	"bufio"
	"math"
	"math/rand/v2"
	"os"
	"strings"
	"testing"
)

func readWords(tb testing.TB) []string {
	tb.Helper()
	f, err := os.Open("testdata/words.txt")
	if err != nil {
		tb.Fatalf("open testdata: %v", err)
	}
	defer f.Close()

	var words []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	if err := scanner.Err(); err != nil {
		tb.Fatalf("read testdata: %v", err)
	}
	return words
}

// TestQualitySample prints a sample of generated words. It is skipped
// unless run with -v. Useful for eyeballing algorithm changes; not a
// pass/fail check.
func TestQualitySample(t *testing.T) {
	t.Parallel()
	if !testing.Verbose() {
		t.Skip("run with -v to see generated samples")
	}

	words := readWords(t)
	dict := &Dictionary{}
	dict.Add(words...)
	g := dict.Generator()
	g.Random = rand.New(rand.NewPCG(1, 0)).Uint32

	const samples = 30
	t.Logf("sample of %d generated words:", samples)
	for range samples {
		t.Logf("  %s", g.WordWithDistance(4, 12))
	}
}

// TestQualityPerplexity reports per-symbol perplexity of the model on
// a held-out 20% split of the fixture. Lower is better; 1.0 means the
// model assigns probability 1 to every test transition (impossible in
// practice). Transitions the model has not seen are charged
// scoreFloor; this floor is a scoring artefact, not a sampler change.
// For an A/B comparison of algorithm variants, run this before and
// after the change and compare the printed numbers.
func TestQualityPerplexity(t *testing.T) {
	t.Parallel()
	words := readWords(t)

	var train, test []string
	for i, w := range words {
		if i%5 == 0 {
			test = append(test, w)
		} else {
			train = append(train, w)
		}
	}

	dict := &Dictionary{}
	dict.Add(train...)
	g := dict.Generator()

	var totalLogP float64
	var totalSyms int
	for _, w := range test {
		word := prefix + strings.ToLower(strings.TrimSpace(w)) + suffix
		logP, n := wordLogProb(g, word)
		totalLogP += logP
		totalSyms += n
	}

	if totalSyms == 0 {
		t.Fatalf("no test words to score")
	}

	avgLogP := totalLogP / float64(totalSyms)
	perplexity := math.Exp(-avgLogP)

	t.Logf("perplexity: %.4f (over %d words, %d symbols, floor %g)",
		perplexity, len(test), totalSyms, scoreFloor)
}

// scoreFloor is the per-symbol probability floor used when a test
// transition has zero probability under the (unsmoothed) model.
// It is a scoring-only patch: it makes perplexity well-defined for
// unseen transitions without changing how Word samples.
const scoreFloor = 1e-3

// wordLogProb returns the natural-log probability of producing word
// (already wrapped with prefix/suffix) under g, using the same backoff
// as Word's sampler. Transitions with zero probability under the model
// are charged scoreFloor instead, so the result is always finite.
func wordLogProb(g Generator, word string) (logP float64, symbols int) {
	maxSeq := g.MaxSequences
	if maxSeq == 0 {
		maxSeq = MaxSequencesDefault
	}

	for i := 1; i < len(word); i++ {
		start := max(i-maxSeq, 0)
		outcomes := g.contextOutcomes([]byte(word[start:i]))
		sym := word[i]

		var p float32
		for j, oc := range outcomes {
			if oc.sym == sym {
				var prev float32
				if j > 0 {
					prev = outcomes[j-1].cum
				}
				p = oc.cum - prev
				break
			}
		}
		if p == 0 {
			p = scoreFloor
		}
		logP += math.Log(float64(p))
		symbols++
	}
	return logP, symbols
}
