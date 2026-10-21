//go:build llgo && (linux || darwin)

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package os_test

import (
	"testing"

	"github.com/goplus/lib/c"
	llos "github.com/goplus/lib/c/os"
)

func TestNativeEnviron(t *testing.T) {
	const key = "LLGO_LIB_NATIVE_ENVIRON_TEST"
	name := c.AllocaCStr(key)
	previous := llos.Getenv(name)
	if previous != nil {
		saved := c.GoString(previous)
		defer func() { llos.Setenv(name, c.AllocaCStr(saved), 1) }()
	} else {
		defer llos.Unsetenv(name)
	}
	if llos.Setenv(name, c.AllocaCStr("visible"), 1) != 0 {
		t.Fatal("setenv failed")
	}
	if llos.Environ == nil {
		t.Fatal("native environ variable is nil")
	}
	for i := 0; ; i++ {
		entry := c.Index(llos.Environ, i)
		if entry == nil {
			t.Fatal("environ did not observe libc setenv")
		}
		if c.GoString(entry) == key+"=visible" {
			return
		}
	}
}
