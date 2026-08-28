//go:build llgo && windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package c_test

import (
	"testing"

	"github.com/goplus/lib/c"
)

func TestWindowsStandardStreams(t *testing.T) {
	if c.Stdin == nil || c.Stdout == nil || c.Stderr == nil {
		t.Fatalf("standard streams are not initialized: stdin=%p stdout=%p stderr=%p", c.Stdin, c.Stdout, c.Stderr)
	}
	if result := c.Fputs(c.AllocaCStr("Windows stdio binding test\n"), c.Stdout); result < 0 {
		t.Fatalf("fputs returned %d", result)
	}
}
