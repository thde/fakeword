package fakeword_test

import (
	"fmt"

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
