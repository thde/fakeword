# fakeword

[![Go Reference](https://pkg.go.dev/badge/github.com/thde/fakeword.svg)](https://pkg.go.dev/github.com/thde/fakeword) [![test](https://github.com/thde/fakeword/actions/workflows/test.yml/badge.svg)](https://github.com/thde/fakeword/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/thde/fakeword)](https://goreportcard.com/report/github.com/thde/fakeword)

Go package fakeword allows to generate fake words.

## Example

Adding some English words, will generate fake words that sound english.

```go
package main

func main() {
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
```

The library was inspired by [nwtgck/go-fakelish](https://github.com/nwtgck/go-fakelish).
