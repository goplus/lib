//go:build llgo && windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package time_test

import (
	"testing"

	"github.com/goplus/lib/c"
	lltime "github.com/goplus/lib/c/time"
)

func TestWindowsTime(t *testing.T) {
	now := lltime.Time(nil)
	if now <= 0 {
		t.Fatalf("_time64 returned %d", now)
	}
	if got := lltime.Difftime(now, now-1); got != 1 {
		t.Fatalf("_difftime64 returned %v, want 1", got)
	}
	if value := lltime.Gmtime(&now); value == nil {
		t.Fatal("_gmtime64 returned nil")
	}
	if value := lltime.Localtime(&now); value == nil {
		t.Fatal("_localtime64 returned nil")
	}
	if value := lltime.Ctime(&now); value == nil || c.GoString(value) == "" {
		t.Fatal("_ctime64 returned an empty C string")
	}
}
