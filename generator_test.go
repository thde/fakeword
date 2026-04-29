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
