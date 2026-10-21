//go:build llgo && (linux || darwin)

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c_test

import (
	"testing"

	"github.com/goplus/lib/c"
)

func TestNativeGlobals(t *testing.T) {
	if c.Stdin == nil || c.Stdout == nil || c.Stderr == nil {
		t.Fatal("native standard stream variables are nil")
	}
	if c.Argc < 1 || c.Argv == nil || c.Index(c.Argv, 0) == nil {
		t.Fatal("native argument variables are not initialized")
	}
	if c.GoString(c.Index(c.Argv, 0)) == "" || c.Index(c.Argv, c.Argc) != nil {
		t.Fatal("native argv is empty or lacks its terminating null")
	}
	previous := c.Opterr
	defer func() { c.Opterr = previous }()
	c.Opterr = 0
	if c.Opterr != 0 {
		t.Fatal("native variable write did not persist")
	}
}
