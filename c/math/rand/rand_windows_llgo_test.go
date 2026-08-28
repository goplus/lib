//go:build llgo && windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package rand_test

import (
	"testing"

	"github.com/goplus/lib/c/math/rand"
)

func TestWindowsRand(t *testing.T) {
	rand.Srand(1)
	first := rand.Rand()
	rand.Srand(1)
	if got := rand.Rand(); got != first {
		t.Fatalf("rand after equal seeds = %d, want %d", got, first)
	}
}
