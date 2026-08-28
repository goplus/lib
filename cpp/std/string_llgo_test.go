//go:build llgo

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package std_test

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/goplus/lib/cpp/std"
)

func TestStringRoundTrip(t *testing.T) {
	if runtime.GOOS == "windows" {
		var stringValue std.String
		if got, want := unsafe.Sizeof(stringValue), uintptr(16)+2*unsafe.Sizeof(uintptr(0)); got != want {
			t.Fatalf("sizeof(std::string) = %d, want %d", got, want)
		}
		if got, want := unsafe.Alignof(stringValue), unsafe.Alignof(uintptr(0)); got != want {
			t.Fatalf("alignof(std::string) = %d, want %d", got, want)
		}
	}

	const value = "LLGo C++ string"
	stringValue := std.NewString(value)
	if got := stringValue.Str(); got != value {
		t.Fatalf("String.Str() = %q, want %q", got, value)
	}
}
