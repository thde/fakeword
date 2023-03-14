# fakeword

[![Go Reference](https://pkg.go.dev/badge/thde.io/fakeword.svg)](https://pkg.go.dev/thde.io/fakeword) [![test](https://github.com/thde/fakeword/actions/workflows/test.yml/badge.svg)](https://github.com/thde/fakeword/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/thde.io/fakeword)](https://goreportcard.com/report/thde.io/fakeword)

Go package fakeword allows to generate fake words.

## Example

Adding some English words, will generate fake words that sound english.

Try the example online on [pkg.go.dev](https://pkg.go.dev/thde.io/fakeword#example-Generator.Word).
```go
package main

import "thde.io/fakeword"

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
