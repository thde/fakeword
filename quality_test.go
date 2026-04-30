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
	for i := 0; i < samples; i++ {
		t.Logf("  %s", g.WordWithDistance(4, 12))
	}
}

// TestQualityPerplexity reports per-symbol perplexity of the model on
// a held-out 20% split of the fixture. Lower is better; 1.0 means the
// model assigns probability 1 to every test transition (impossible in
// practice). For an A/B comparison of algorithm variants, run this
// before and after the change and compare the printed numbers.
func TestQualityPerplexity(t *testing.T) {
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
	skipped := 0
	for _, w := range test {
		word := prefix + strings.ToLower(strings.TrimSpace(w)) + suffix
		ok, logP, n := wordLogProb(g, word)
		if !ok {
			skipped++
			continue
		}
		totalLogP += logP
		totalSyms += n
	}

	if totalSyms == 0 {
		t.Fatalf("no test words could be scored (all had unseen contexts)")
	}

	avgLogP := totalLogP / float64(totalSyms)
	perplexity := math.Exp(-avgLogP)

	t.Logf("perplexity: %.4f (scored %d/%d words, %d symbols, %d skipped due to zero probability)",
		perplexity, len(test)-skipped, len(test), totalSyms, skipped)
}

// wordLogProb returns the natural-log probability of producing word
// (already wrapped with prefix/suffix) under g, using the same backoff
// as Word's sampler. ok is false if any transition has probability 0.
func wordLogProb(g Generator, word string) (ok bool, logP float64, symbols int) {
	maxSeq := g.MaxSequences
	if maxSeq == 0 {
		maxSeq = MaxSequencesDefault
	}

	for i := 1; i < len(word); i++ {
		start := i - maxSeq
		if start < 0 {
			start = 0
		}
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
			return false, 0, 0
		}
		logP += math.Log(float64(p))
		symbols++
	}
	return true, logP, symbols
}
