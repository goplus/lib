//go:build llgo && windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package setjmp_test

import (
	"testing"

	"github.com/goplus/lib/c/setjmp"
)

func TestWindowsSetjmp(t *testing.T) {
	var env setjmp.JmpBuf
	switch result := setjmp.Setjmp(&env); result {
	case 0:
		setjmp.Longjmp(&env, 7)
		t.Fatal("longjmp returned")
	case 7:
		return
	default:
		t.Fatalf("setjmp returned %d, want 0 then 7", result)
	}
}
