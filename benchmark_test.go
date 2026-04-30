package fakeword

import (
	"math/rand/v2"
	"os"
	"sync/atomic"
	"testing"
)

func loadDict(tb testing.TB) *Dictionary {
	tb.Helper()
	f, err := os.Open("testdata/words.txt")
	if err != nil {
		tb.Fatalf("open testdata: %v", err)
	}
	tb.Cleanup(func() { _ = f.Close() })
	return (&Dictionary{}).Read(f)
}

func BenchmarkWord(b *testing.B) {
	g := loadDict(b).Generator()
	g.Random = rand.New(rand.NewPCG(1, 0)).Uint32
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Word()
	}
}

func BenchmarkWordWithDistance(b *testing.B) {
	g := loadDict(b).Generator()
	g.Random = rand.New(rand.NewPCG(1, 0)).Uint32
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.WordWithDistance(6, 10)
	}
}

func BenchmarkWord_ManualProbs(b *testing.B) {
	src := loadDict(b).Generator()
	g := Generator{
		Probabilities: src.Probabilities,
		MaxSequences:  src.MaxSequences,
		Random:        rand.New(rand.NewPCG(1, 0)).Uint32,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Word()
	}
}

func BenchmarkDictionary_Generator(b *testing.B) {
	dict := loadDict(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dict.Generator()
	}
}

// BenchmarkWord_Parallel shares a single Generator across goroutines.
func BenchmarkWord_Parallel(b *testing.B) {
	g := loadDict(b).Generator()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = g.Word()
		}
	})
}

// BenchmarkWord_Parallel_PerGoroutineRand gives each goroutine its own
// seeded Random func.
func BenchmarkWord_Parallel_PerGoroutineRand(b *testing.B) {
	base := loadDict(b).Generator()
	var seed int64
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		g := base
		g.Random = rand.New(rand.NewPCG(uint64(atomic.AddInt64(&seed, 1)), 0)).Uint32
		for pb.Next() {
			_ = g.Word()
		}
	})
}
