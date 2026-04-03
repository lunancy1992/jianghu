package crawler

import (
	"hash/fnv"
	"math/bits"
	"strings"
	"unicode"
)

// SimHash computes a 64-bit SimHash of the given text.
func SimHash(text string) uint64 {
	tokens := tokenize(text)
	var v [64]int

	for _, token := range tokens {
		h := hashToken(token)
		for i := 0; i < 64; i++ {
			if (h>>uint(i))&1 == 1 {
				v[i]++
			} else {
				v[i]--
			}
		}
	}

	var fingerprint uint64
	for i := 0; i < 64; i++ {
		if v[i] > 0 {
			fingerprint |= 1 << uint(i)
		}
	}
	return fingerprint
}

// HammingDistance returns the number of differing bits between two hashes.
func HammingDistance(a, b uint64) int {
	return bits.OnesCount64(a ^ b)
}

// IsDuplicate returns true if the hamming distance is less than threshold.
func IsDuplicate(a, b uint64, threshold int) bool {
	return HammingDistance(a, b) < threshold
}

func tokenize(text string) []string {
	var tokens []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	// Generate bigrams for better CJK handling
	n := len(tokens)
	if n > 1 {
		for i := 0; i < n-1; i++ {
			tokens = append(tokens, tokens[i]+tokens[i+1])
		}
	}

	return tokens
}

func hashToken(token string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(token))
	return h.Sum64()
}
