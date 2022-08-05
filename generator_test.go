package fakeword

import (
	"math/rand"
	"testing"
)

func TestGenerator_Word(t *testing.T) {
	// ensure we get reproducable results
	random = rand.New(rand.NewSource(1))

	type fields struct {
		Probabilities map[string]map[string]float32
		MaxSequences  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple",
			fields: fields{
				Probabilities: map[string]map[string]float32{"^": {"b": 1.0}, "b": {"$": 1.0}},
				MaxSequences: 2,
			},
			want: "b",
		},
		{
			name: "empty",
			fields: fields{
				Probabilities: map[string]map[string]float32{},
				MaxSequences: 2,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Generator{
				Probabilities: tt.fields.Probabilities,
				MaxSequences:  tt.fields.MaxSequences,
			}
			if got := g.Word(); got != tt.want {
				t.Errorf("Generator.Word() = %v, want %v", got, tt.want)
			}
		})
	}
}
