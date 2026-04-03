package crawler

import (
	"testing"
)

func TestSimHash_SameText(t *testing.T) {
	text := "今日武林大会在少林寺召开，各大门派掌门齐聚一堂"
	h1 := SimHash(text)
	h2 := SimHash(text)
	if h1 != h2 {
		t.Errorf("SimHash of same text should be equal: %d != %d", h1, h2)
	}
}

func TestSimHash_SimilarText(t *testing.T) {
	// Same text with minor suffix addition should produce identical hash
	t1 := "the quick brown fox jumps over the lazy dog near the river bank"
	t2 := "the quick brown fox jumps over the lazy dog near the river bank today"
	h1 := SimHash(t1)
	h2 := SimHash(t2)
	dist := HammingDistance(h1, h2)
	if dist > 15 {
		t.Errorf("Similar texts should have small hamming distance, got %d", dist)
	}
}

func TestSimHash_DifferentText(t *testing.T) {
	t1 := "今日武林大会在少林寺召开"
	t2 := "The quick brown fox jumps over the lazy dog"
	h1 := SimHash(t1)
	h2 := SimHash(t2)
	dist := HammingDistance(h1, h2)
	if dist < 5 {
		t.Errorf("Very different texts should have large hamming distance, got %d", dist)
	}
}

func TestHammingDistance(t *testing.T) {
	tests := []struct {
		a, b uint64
		want int
	}{
		{0, 0, 0},
		{0, 1, 1},
		{0b1111, 0b0000, 4},
		{0xFFFFFFFFFFFFFFFF, 0, 64},
	}
	for _, tt := range tests {
		got := HammingDistance(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("HammingDistance(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestIsDuplicate(t *testing.T) {
	if !IsDuplicate(0, 1, 3) {
		t.Error("IsDuplicate(0, 1, 3) should be true (distance 1 < 3)")
	}
	if IsDuplicate(0, 0xFFFFFFFFFFFFFFFF, 3) {
		t.Error("IsDuplicate(0, max, 3) should be false (distance 64 >= 3)")
	}
}

func TestTokenize(t *testing.T) {
	tokens := tokenize("hello world")
	if len(tokens) < 2 {
		t.Fatalf("tokenize should return at least 2 tokens, got %d", len(tokens))
	}
	if tokens[0] != "hello" || tokens[1] != "world" {
		t.Errorf("tokenize('hello world') = %v, want [hello, world, ...]", tokens)
	}
	// Should also have bigram
	found := false
	for _, tok := range tokens {
		if tok == "helloworld" {
			found = true
			break
		}
	}
	if !found {
		t.Error("tokenize should generate bigrams")
	}
}

func TestTokenize_Empty(t *testing.T) {
	tokens := tokenize("")
	if len(tokens) != 0 {
		t.Errorf("tokenize('') should return empty, got %v", tokens)
	}
}

func TestSimHash_EmptyText(t *testing.T) {
	h := SimHash("")
	if h != 0 {
		t.Errorf("SimHash('') should be 0, got %d", h)
	}
}
