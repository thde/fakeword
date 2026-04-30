package fakeword_test

import (
	"fmt"
	"sync"

	"thde.io/fakeword"
)

func ExampleGenerator_Word() {
	words := []string{
		"Psychotomimetic",
		"Pulchritudinous",
		"Consanguineous",
		"Trichotillomania",
	}

	dict := fakeword.Dictionary{}
	dict.Add(words...)

	gen := dict.Generator()
	fmt.Println(gen.Word())
}

// ExampleGenerator_Word_concurrent shows how to generate words from
// multiple goroutines. The default RNG is safe for concurrent use, so
// sharing a Generator across goroutines requires no extra setup.
func ExampleGenerator_Word_concurrent() {
	words := []string{
		"Psychotomimetic",
		"Pulchritudinous",
		"Consanguineous",
		"Trichotillomania",
	}

	dict := fakeword.Dictionary{}
	dict.Add(words...)
	base := dict.Generator()

	const workers = 4
	var wg sync.WaitGroup
	results := make(chan string, workers)

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- base.Word()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var generated []string
	for w := range results {
		generated = append(generated, w)
	}
	fmt.Printf("generated %d fake words\n", len(generated))
	// Output: generated 4 fake words
}
