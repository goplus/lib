//go:build llgo && !windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package time_test

import (
	"testing"

	"github.com/goplus/lib/c"
	lltime "github.com/goplus/lib/c/time"
)

func TestCtime(t *testing.T) {
	now := lltime.Time(nil)
	if value := lltime.Ctime(&now); value == nil || c.GoString(value) == "" {
		t.Fatal("ctime returned an empty C string")
	}
}
