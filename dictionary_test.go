package fakeword

import (
	"reflect"
	"testing"
)

func TestDictionary_Add(t *testing.T) {
	tests := []struct {
		name         string
		prefixLength int
		words        []string
		want         *Dictionary
	}{
		{
			"simple",
			4,
			[]string{"hi"},
			&Dictionary{
				counter: map[string]map[string]int{
					"^": {"h": 1},
					"^h": {"i": 1},
					"^hi": {"$": 1},
					"h": {"i": 1},
					"hi": {"$": 1},
					"i": {"$": 1},
				},
				PrefixLength: 4,
			},
		},
		{
			"short prefix",
			2,
			[]string{"hi"},
			&Dictionary{
				counter: map[string]map[string]int{
					"^": {"h": 1},
					"h": {"i": 1},
					"i": {"$": 1},
				},
				PrefixLength: 2,
			},
		},
		{
			"empty",
			defaultPrefixLength,
			[]string{},
			&Dictionary{
				PrefixLength: defaultPrefixLength,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Dictionary{PrefixLength: tt.prefixLength}
			if got := w.Add(tt.words...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dictionary.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictionary_Generator(t *testing.T) {
	tests := []struct {
		name    string
		counter map[string]map[string]int
		want    map[string]map[string]float32
	}{
		{
			"simple",
			map[string]map[string]int{"a": {"b": 2, "c": 2}},
			map[string]map[string]float32{"a": {"b": 0.5, "c": 0.5}},
		},
		{
			"simple",
			map[string]map[string]int{"a": {"b": 2, "c": 2}, "b": {"c": 2}},
			map[string]map[string]float32{"a": {"b": 0.5, "c": 0.5}, "b": {"c": 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Dictionary{
				counter: tt.counter,
			}

			got := w.Generator()
			if !reflect.DeepEqual(got, Generator{Probabilities: tt.want}) {
				t.Errorf("Words.Generator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictionary_count(t *testing.T) {
	tests := []struct {
		name    string
		substrs []string
		want    map[string]map[string]int
	}{
		{
			"simple",
			[]string{"ab", "ab"},
			map[string]map[string]int{"a": {"b": 2}},
		},
		{
			"simple",
			[]string{"ab", "ab", "ac", "bc"},
			map[string]map[string]int{"a": {"b": 2, "c": 1}, "b": {"c": 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Dictionary{}
			for _, substr := range tt.substrs {
				w.count(substr)
			}
			if got := w.counter; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Words.count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitToLength(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"basic", args{"test", 2}, []string{"te", "es", "st"}},
		{"empty", args{"", 2}, []string{}},
		{"single", args{"test", 1}, []string{"t", "e", "s", "t"}},
		{"equal length", args{"test", 4}, []string{"test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitToLength(tt.args.s, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitToLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
