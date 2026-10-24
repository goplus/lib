//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package signal

import (
	"testing"

	"github.com/goplus/lib/c"
)

// The C signal callback only stores the signal number; it does not allocate.
var deliveredSignal c.Int

type signalABISlot struct {
	size, align uintptr
	class       int
}

func TestNativeSignalConstants(t *testing.T) {
	if SIGUSR1 != nativeUserSignal {
		t.Fatalf("SIGUSR1 = %d, want %d", SIGUSR1, nativeUserSignal)
	}
	if SIGEV_NONE != nativeNoNotification {
		t.Fatalf("SIGEV_NONE = %d, want %d", SIGEV_NONE, nativeNoNotification)
	}
}

func TestNativeSignalSetOperations(t *testing.T) {
	var set SigsetT
	if result := nativeEmptySet(&set); result != 0 {
		t.Fatalf("sigemptyset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 0 {
		t.Fatalf("empty membership = %d, want 0", got)
	}
	if result := nativeAddSet(&set, nativeUserSignal); result != 0 {
		t.Fatalf("sigaddset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 1 {
		t.Fatalf("added membership = %d, want 1", got)
	}
	if result := nativeDeleteSet(&set, nativeUserSignal); result != 0 {
		t.Fatalf("sigdelset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 0 {
		t.Fatalf("deleted membership = %d, want 0", got)
	}
	if result := nativeFillSet(&set); result != 0 {
		t.Fatalf("sigfillset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 1 {
		t.Fatalf("filled membership = %d, want 1", got)
	}
}
