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
