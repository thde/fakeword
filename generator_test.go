package fakeword

import (
	"math/rand/v2"
	"testing"
)

func TestGenerator_Word(t *testing.T) {
	rand := rand.New(rand.NewPCG(1, 2))

	tests := []struct {
		name      string
		generator Generator
		want      string
	}{
		{
			name: "simple",
			generator: Generator{
				Probabilities: map[string]map[string]float32{"^": {"b": 1.0}, "b": {"$": 1.0}},
				MaxSequences:  2,
			},
			want: "b",
		},
		{
			name: "empty",
			generator: Generator{
				Probabilities: map[string]map[string]float32{},
				MaxSequences:  2,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.generator.Random = rand.Uint32

			if got := tt.generator.Word(); got != tt.want {
				t.Errorf("Generator.Word() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_WordWithDistance(t *testing.T) {
	probs := map[string]map[string]float32{
		"^": {"a": 0.5, "b": 0.5},
		"a": {"b": 0.4, "c": 0.4, "$": 0.2},
		"b": {"a": 0.4, "c": 0.4, "$": 0.2},
		"c": {"a": 0.3, "b": 0.3, "$": 0.4},
	}

	tests := []struct {
		name     string
		min, max int
	}{
		{"narrow band", 3, 5},
		{"single length", 4, 4},
		{"min only", 5, 100},
		{"min zero", 0, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Generator{
				Probabilities: probs,
				MaxSequences:  2,
				Random:        rand.New(rand.NewPCG(1, 0)).Uint32,
			}
			for i := 0; i < 200; i++ {
				w := g.WordWithDistance(tt.min, tt.max)
				if len(w) < tt.min || len(w) > tt.max {
					t.Errorf("got len(%q)=%d, want in [%d, %d]", w, len(w), tt.min, tt.max)
				}
			}
		})
	}
}

func TestGenerator_Word_Reproducible(t *testing.T) {
	probs := map[string]map[string]float32{
		"^": {"a": 0.4, "b": 0.3, "c": 0.3},
		"a": {"b": 0.5, "c": 0.5},
		"b": {"a": 0.5, "c": 0.5},
		"c": {"$": 1.0},
	}

	var first string
	for i := 0; i < 5; i++ {
		g := Generator{
			Probabilities: probs,
			MaxSequences:  2,
			Random:        rand.New(rand.NewPCG(42, 0)).Uint32,
		}
		got := g.Word()
		if i == 0 {
			first = got
			continue
		}
		if got != first {
			t.Errorf("seeded Word() not reproducible: run 0 = %q, run %d = %q", first, i, got)
		}
	}
}
