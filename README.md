# fakeword

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
